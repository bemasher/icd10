package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bemasher/icd10/util"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
)

const (
	ResultLimit = 250

	defaultQuerySet = "cm"
)

var (
	docDb *bolt.DB

	tmpls map[string]*template.Template

	production bool
	hostnames  string
	certDir    string

	commit string // set to git rev-parse --short HEAD
)

func SearchIdx(idxBkt, qry string) (docIds []string, err error) {
	tx, err := docDb.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "docDb.Begin")
	}
	defer tx.Commit()

	bkt := tx.Bucket([]byte(idxBkt))
	var matches util.DocIDMap

	for _, token := range util.Tokenize(qry) {
		tokenDocIds := util.DocIDMap{}

		if token.IsPrefix {
			for _, form := range token.Forms {
				docIds, err := SearchPrefix(bkt, form)
				if err != nil {
					return nil, errors.Wrap(err, "SearchPrefix")
				}
				for docId := range docIds {
					tokenDocIds[docId] = true
				}
			}
		} else {
			for _, form := range token.Forms {
				docIds, err := SearchTerm(bkt, form)
				if err != nil {
					return nil, errors.Wrap(err, "SearchTerm")
				}
				for docId := range docIds {
					tokenDocIds[docId] = true
				}
			}
		}

		if matches == nil {
			matches = tokenDocIds
			continue
		}

		if token.IsNegative {
			for docId := range tokenDocIds {
				delete(matches, docId)
			}
		} else {
			for docId := range matches {
				if !tokenDocIds[docId] {
					delete(matches, docId)
				}
			}
		}
	}

	for docId := range matches {
		docIds = append(docIds, docId)
	}

	return
}

func SearchTerm(bkt *bolt.Bucket, term string) (docIds util.DocIDMap, err error) {
	v := bkt.Get([]byte(term))
	if v == nil {
		return nil, nil
	}

	_, err = docIds.UnmarshalMsg(v)
	return docIds, errors.Wrap(err, "docIds.UnmarshalMsg")
}

func SearchPrefix(bkt *bolt.Bucket, prefix string) (docIds util.DocIDMap, err error) {
	docIds = make(util.DocIDMap)

	c := bkt.Cursor()
	p := []byte(prefix)
	for k, v := c.Seek(p); k != nil && bytes.HasPrefix(k, p); k, v = c.Next() {
		var prefixDocIds util.DocIDMap
		_, err = prefixDocIds.UnmarshalMsg(v)
		if err != nil {
			return nil, errors.Wrap(err, "prefixDocIds.UnmarshalMsg")
		}

		for k := range prefixDocIds {
			docIds[k] = true
		}
	}

	return
}

func SearchDocs(docBkt, idxBkt, qry string, unmarshal func([]byte) error) (err error) {
	tx, err := docDb.Begin(false)
	if err != nil {
		return errors.Wrap(err, "docDb.Begin")
	}
	defer tx.Commit()

	docs := tx.Bucket([]byte(docBkt))

	docIds, err := SearchIdx(idxBkt, qry)
	if err != nil {
		return errors.Wrap(err, "SearchIdx")
	}
	for _, docId := range docIds {
		err = unmarshal(docs.Get([]byte(docId)))
		if err != nil {
			return errors.Wrap(err, "unmarshal")
		}
	}

	return nil
}

type TemplateData struct {
	*sync.Mutex

	TmplName string
	Commit   string

	Query   string
	Stemmed string

	Results map[string]interface{}

	ResultLimit int
	Duration    time.Duration
	Errors      []error
}

type IndexHandler struct{}

func (IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	start := time.Now()

	tmplData := TemplateData{
		Mutex:       new(sync.Mutex),
		TmplName:    "index",
		Commit:      commit,
		Results:     map[string]interface{}{},
		ResultLimit: ResultLimit,
	}

	// If the commit value isn't set, supply a value that should always cache-bust.
	if tmplData.Commit == "" {
		tmplData.Commit = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	// If we're in development mode, parse the templates for each request.
	if !production {
		err = ParseTemplates()
		if err != nil {
			log.Fatalf("%+v\n", errors.Wrap(err, "ParseTemplates"))
		}
	}

	query := r.URL.Query()

	if q, ok := query["q"]; ok {
		tmplData.Query = q[0]
		tmplData.Stemmed = fmt.Sprintf("%+v\n", util.Tokenize(q[0]))
	}

	if queryFns, ok := querySets[defaultQuerySet]; ok {
		tmplData.TmplName = defaultQuerySet

		wg := new(sync.WaitGroup)
		wg.Add(len(queryFns))

		for _, queryFn := range queryFns {
			go func(q Query) {
				tmplData.Lock()
				tmplData.Results[q.Name], err = q.Fn(tmplData.Query)
				tmplData.Unlock()
				if err != nil {
					tmplData.Errors = append(tmplData.Errors, errors.Wrap(err, "q.Fn"))
				}

				wg.Done()
			}(queryFn)
		}

		wg.Wait()
	}

	tmplData.Duration = time.Now().Sub(start).Round(time.Microsecond)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpls[tmplData.TmplName].ExecuteTemplate(w, "index.html", tmplData)
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "templates.ExecuteTemplate"))
	}
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL.String())

		h.ServeHTTP(w, r)
	})
}

func CacheControl(h http.Handler, val string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", val)

		h.ServeHTTP(w, r)
	})
}

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gzipWriter := gzipWriterPool.Get().(*gzip.Writer)
		defer gzipWriterPool.Put(gzipWriter)

		gzipWriter.Reset(w)
		defer gzipWriter.Close()

		next.ServeHTTP(&gzipResponseWriter{gzipWriter, w}, r)
	})
}

func Source(src string) template.HTML {
	switch src {
	case "alpha":
		return template.HTML("Alpha")
	case "drug":
		return template.HTML("Drug")
	case "ext":
		return template.HTML("External")
	case "neo":
		return template.HTML("Neoplasm")
	}
	return ""
}

func Badge(kind string) template.HTML {
	switch kind {
	case "notes":
		return template.HTML("Notes")
	case "includes":
		return template.HTML("Includes")
	case "inclusionTerm":
		return template.HTML("Inclusion Term")
	case "codeFirst":
		return template.HTML("Code First")
	case "codeAlso":
		return template.HTML("Code Also")
	case "useAdditionalCode":
		return template.HTML("Use Additional")
	case "excludes1":
		return template.HTML("Excludes1")
	case "excludes2":
		return template.HTML("Excludes2")
	case "sevenChrNote":
		return template.HTML("7<sup>th</sup> Character")

	case "see":
		return template.HTML("See")
	case "seeAlso":
		return template.HTML("See Also")
	case "seecat":
		return template.HTML("See Category")
	case "subcat":
		return template.HTML("Sub Category")
	}

	return ""
}

func Label(kind string) template.HTML {
	switch kind {
	case "notes":
		return template.HTML("<b>Notes:</b></br>")
	case "includes":
		return template.HTML("<b>Includes:</b></br>")
	case "inclusionTerm":
		return template.HTML("<b>Inclusion Term:</b></br>")
	case "codeFirst":
		return template.HTML("<b>Code First:</b></br>")
	case "codeAlso":
		return template.HTML("<b>Code Also:</b></br>")
	case "useAdditionalCode":
		return template.HTML("<b>Use Additional</b>")
	case "excludes1":
		return template.HTML("<b>Excludes1:</b></br>")
	case "excludes2":
		return template.HTML("<b>Excludes2:</b></br>")
	case "sevenChrNote":
		return template.HTML("<b>7<sup>th</sup> Character:</b>")

	case "see":
		return template.HTML("<b>See:</b>&nbsp;")
	case "seeAlso":
		return template.HTML("<b>See Also:</b>&nbsp;")
	case "seecat":
		return template.HTML("<b>See Category:</b>&nbsp;")
	case "subcat":
		return template.HTML("<b>Sub Category:</b>&nbsp;")
	}

	return ""
}

func CodeTrimSuffix(code string) string {
	return strings.TrimRight(code, ".-")
}

func ParseTemplates() (err error) {
	tmpls = map[string]*template.Template{}

	tmpl := template.New("").Funcs(template.FuncMap{
		"badge":    Badge,
		"label":    Label,
		"source":   Source,
		"codeTrim": CodeTrimSuffix,
	})

	tmpl, err = tmpl.ParseFiles("tmpl/index.html")
	if err != nil {
		return errors.Wrap(err, "tmpl.ParseFiles")
	}
	tmpls["index"] = tmpl

	for _, t := range []struct{ name, filename string }{
		{"cm", "tmpl/cm.html"},
	} {
		tmpls[t.name], err = tmpl.Clone()
		if err != nil {
			return errors.Wrap(err, "indexTmpl.Clone")
		}
		tmpls[t.name], err = tmpls[t.name].ParseFiles(t.filename)
		if err != nil {
			return errors.Wrap(err, "tmpls[t.name].ParseFiles")
		}
	}

	return nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	flag.BoolVar(&production, "prod", false, "run with production settings")
	flag.StringVar(
		&hostnames,
		"hosts",
		"example.com",
		"comma-separated list of hostnames to obtain tls certificates for",
	)

	flag.StringVar(&certDir, "certdir", ".cert", "directory to store tls certs")
	flag.Parse()

	var err error

	err = ParseTemplates()
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "ParseTemplates"))
	}

	docDb, err = bolt.Open("documents.db", 0600, nil)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "bolt.Open"))
	}
}

func main() {
	defer docDb.Close()

	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(strings.Split(hostnames, ",")...),
	}

	server := http.Server{
		ErrorLog: log.New(os.Stderr, "http: ", log.Lshortfile|log.Lmicroseconds),
	}

	mux := http.NewServeMux()
	mux.Handle("/", Logger(Gzip(IndexHandler{})))
	mux.Handle("/assets/", Gzip(
		CacheControl(
			http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))),
			"public, max-age=31536000",
		),
	))

	if production {
		server.Addr = ":https"
		server.Handler = m.HTTPHandler(mux)
		server.TLSConfig = m.TLSConfig()

		go func() {
			log.Fatal(http.ListenAndServe(":http", m.HTTPHandler(nil)))
		}()

		log.Printf("Production Listening on: %s\n", server.Addr)
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("%+v\n", errors.Wrap(err, "server.ListenAndServeTLS"))
		}
	} else {
		server.Addr = ":8080"
		server.Handler = mux

		log.Printf("Development Listening on: %s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("%+v\n", errors.Wrap(err, "server.ListenAndServe"))
		}
	}
}
