package main

import (
	"compress/gzip"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
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
)

var (
	docDb     *bolt.DB
	templates *template.Template
)

func Search(bucketId, qry string) (docIds []string, err error) {
	tx, err := docDb.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "docDb.Begin")
	}
	defer tx.Commit()

	bkt := tx.Bucket([]byte(bucketId))
	docFreq := map[string]int{}

	terms := util.Tokenize(qry)
	for _, term := range terms {
		term = strings.TrimRight(term, ".")

		b := bkt.Get([]byte(term))
		if len(b) == 0 {
			continue
		}

		var termDocIds util.DocIDMap
		_, err = termDocIds.UnmarshalMsg(b)
		if err != nil {
			return nil, errors.Wrap(err, "termDocIds.UnmarshalMsg")
		}

		for docId := range termDocIds {
			docFreq[docId]++
		}
	}

	for docId, freq := range docFreq {
		if freq == len(terms) {
			docIds = append(docIds, docId)
		}
	}

	return
}

func AlphabeticQuery(qry string) (terms []util.Term, err error) {
	tx, err := docDb.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "docDb.Begin")
	}
	defer tx.Commit()

	docs := tx.Bucket([]byte("alphabetic"))

	docIds, err := Search("alphabetic_index", qry)
	if err != nil {
		return nil, errors.Wrap(err, "Search")
	}
	for _, docId := range docIds {
		var term util.Term
		_, err = term.UnmarshalMsg(docs.Get([]byte(docId)))
		if err != nil {
			return nil, errors.Wrap(err, "term.UnmarshalMsg")
		}

		terms = append(terms, term)
	}

	sort.Slice(terms, func(i, j int) bool {
		return strings.Compare(terms[i].Title, terms[j].Title) < 0
	})

	trim := ResultLimit
	if len(terms) < trim {
		trim = len(terms)
	}

	return terms[:trim], nil
}

func TabularQuery(qry string) (diags []util.Diag, err error) {
	tx, err := docDb.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "docDb.Begin")
	}
	defer tx.Commit()

	docs := tx.Bucket([]byte("tabular"))

	docIds, err := Search("tabular_index", qry)
	if err != nil {
		return nil, errors.Wrap(err, "Search")
	}
	for _, docId := range docIds {
		var diag util.Diag
		_, err = diag.UnmarshalMsg(docs.Get([]byte(docId)))
		if err != nil {
			return nil, errors.Wrap(err, "term.UnmarshalMsg")
		}

		diags = append(diags, diag)
	}

	sort.Slice(diags, func(i, j int) bool {
		return strings.Compare(diags[i].Code, diags[j].Code) < 0
	})

	trim := ResultLimit
	if len(diags) < trim {
		trim = len(diags)
	}

	return diags[:trim], nil
}

type QueryResults struct {
	Query   string
	Stemmed string

	Alphabetic []util.Term
	Tabular    []util.Diag

	AlphaTrunc bool
	TabTrunc   bool
	Duration   time.Duration
	Errors     []error
}

type IndexHandler struct{}

func (IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	query := r.URL.Query()

	var terms string
	if _, ok := query["q"]; ok {
		terms = query["q"][0]
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	results := QueryResults{
		Query:   terms,
		Stemmed: strings.Join(util.Tokenize(terms), " "),
	}

	go func() {
		var err error
		results.Alphabetic, err = AlphabeticQuery(results.Query)
		results.Errors = append(results.Errors, err)
		results.AlphaTrunc = len(results.Alphabetic) == ResultLimit
		wg.Done()
	}()

	go func() {
		var err error
		results.Tabular, err = TabularQuery(results.Query)
		results.Errors = append(results.Errors, err)
		results.TabTrunc = len(results.Tabular) == ResultLimit
		wg.Done()
	}()

	wg.Wait()
	results.Duration = time.Now().Sub(start).Round(time.Microsecond)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "index.html", results)
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "templates.ExecuteTemplate"))
	}
}

type Adapter func(http.Handler) http.Handler

func Logger() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr, r.Method, r.URL.String())

			h.ServeHTTP(w, r)
		})
	}
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

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	var err error

	templates = template.New("").Funcs(template.FuncMap{
		"badge":    Badge,
		"label":    Label,
		"codeTrim": CodeTrimSuffix,
	})
	templates = template.Must(templates.ParseGlob("assets/*.html"))

	docDb, err = bolt.Open("documents.db", 0600, nil)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "bolt.Open"))
	}
}

func main() {
	production := flag.Bool("prod", false, "run with production settings")
	hostnames := flag.String(
		"hosts",
		"example.com",
		"comma-separated list of hostnames to obtain tls certificates for",
	)
	certDir := flag.String("certdir", ".cert", "directory to store tls certs")
	flag.Parse()

	defer docDb.Close()

	m := &autocert.Manager{
		Cache:      autocert.DirCache(*certDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(strings.Split(*hostnames, ",")...),
	}

	server := http.Server{
		ErrorLog: log.New(os.Stderr, "http: ", log.Lshortfile|log.Lmicroseconds),
	}

	mux := http.NewServeMux()
	mux.Handle("/", Logger()(Gzip(IndexHandler{})))
	mux.Handle("/assets/", Gzip(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))))

	if *production {
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
		server.Addr = "127.0.0.1:8080"
		server.Handler = mux

		log.Printf("Development Listening on: %s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("%+v\n", errors.Wrap(err, "server.ListenAndServe"))
		}
	}
}
