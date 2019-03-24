package main

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"encoding/xml"
	"log"
	"os"
	"strings"

	"github.com/bemasher/icd10/util"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const (
	xmlFilename = "icd10cm_index_2019.xml"

	boltFilename = "../documents.db"
	boltBucket   = "alphabetic"
)

type Node struct {
	Title Title `xml:"title"`

	XMLName  xml.Name
	Content  string `xml:",chardata"`
	Children []Node `xml:",any"`
}

func (n Node) Walk(ctx context.Context) chan util.Term {
	terms := make(chan util.Term)

	go func() {
		defer close(terms)
		for _, child := range n.Children {
			select {
			case <-ctx.Done():
				return
			default:
				child.walk(ctx, terms, 0, nil)
			}
		}
	}()

	return terms
}

func (n Node) walk(ctx context.Context, terms chan util.Term, depth int, breadcrumbs []string) {
	if len(n.Children) == 0 {
		n.Children = nil
	}

	select {
	case <-ctx.Done():
		return
	default:
		switch elem := n.XMLName.Local; elem {
		case "mainTerm", "term":
			breadcrumbs = append(breadcrumbs, n.Title.Title)

			var term util.Term

			term.Title = strings.Join(breadcrumbs, ", ")
			term.Nemod = n.Title.Nemod

			for _, child := range n.Children {
				childName := child.XMLName.Local

				switch childName {
				case "mainTerm", "term":
					continue
				case "code":
					term.Code = child.Content
				case "manif":
					term.Manif = child.Content
				default:
					term.Attrs = append(term.Attrs, util.Attr{childName, child.Content})
				}
			}

			terms <- term
			fallthrough
		default:
			for _, child := range n.Children {
				child.walk(ctx, terms, depth+1, breadcrumbs)
			}
		}
	}
}

type Title struct {
	Title string `xml:",chardata" json:",omitempty"`
	Nemod string `xml:"nemod" json:",omitempty"`
}

func WriteAttr(w *csv.Writer, termId, attr string) {
	if attr == "" {
		return
	}

	w.Write([]string{termId, attr})
}

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	xmlFile, err := os.Open(xmlFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	xmlDecoder := xml.NewDecoder(xmlFile)

	var node Node
	err = xmlDecoder.Decode(&node)
	if err != nil {
		log.Fatal(err)
	}

	docDb, err := bolt.Open(boltFilename, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer docDb.Close()

	docDb.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(boltBucket))
		return nil
	})

	tx, err := docDb.Begin(true)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "docDb.Begin"))
	}
	defer tx.Commit()

	bucket, err := tx.CreateBucket([]byte(boltBucket))
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "tx.CreateBucket"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	idx := 0
	for term := range node.Walk(ctx) {
		docId := make([]byte, 8)
		docIdLen := binary.PutUvarint(docId, uint64(idx))

		doc, err := term.MarshalMsg(nil)
		if err != nil {
			log.Fatalf("%+v\n", errors.Wrap(err, "json.Marshal"))
		}

		bucket.Put(docId[:docIdLen], doc)

		idx++
	}

	log.Println("parsed", idx, "terms")
}
