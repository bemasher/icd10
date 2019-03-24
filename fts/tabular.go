package main

import (
	"context"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/bemasher/icd10/util"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const (
	xmlFilename = "icd10cm_tabular_2019.xml"

	boltFilename = "../documents.db"
	boltBucket   = "tabular"

	batchSize = 1000
)

type Node struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`

	Code string `xml:"name"`
	Desc string `xml:"desc"`

	Notes       []string    `xml:"note"`
	SevenChrDef []Extension `xml:"sevenChrDef>extension"`

	Children []Node `xml:",any"`
}

type Extension struct {
	Char string `xml:"char,attr"`
	Def  string `xml:",innerxml"`
}

func (e Extension) String() string {
	return fmt.Sprintf("%q - %s", e.Char, e.Def)
}

func (n Node) Walk(ctx context.Context) chan util.Diag {
	diags := make(chan util.Diag)

	go func() {
		for _, child := range n.Children {
			child.walk(ctx, diags, 0, nil)
		}
		close(diags)
	}()

	return diags
}

func (n Node) walk(ctx context.Context, diags chan util.Diag, depth int, parent *Node) {
	if len(n.Children) == 0 {
		n.Children = nil
	}

	select {
	case <-ctx.Done():
		return
	default:
		switch elem := n.XMLName.Local; elem {
		case "version":
			return
		case "introduction", "sectionIndex":
			return
		case "visCategory", "visMax", "visMin", "visRange", "visualImpairment":
			return
		case "diag":
			// If this node doesn't have any 7th character definitions, and the
			// parent does, copy them.
			if len(n.SevenChrDef) == 0 && parent != nil {
				n.SevenChrDef = make([]Extension, len(parent.SevenChrDef))
				copy(n.SevenChrDef, parent.SevenChrDef)

				for _, child := range parent.Children {
					if child.XMLName.Local == "sevenChrNote" {
						n.Children = append(n.Children, child)
					}
				}
			}

			diag := util.Diag{
				Code: n.Code,
				Desc: n.Desc,
			}

			for _, child := range n.Children {
				if len(child.Notes) != 0 {
					note := util.Note{Kind: child.XMLName.Local}
					note.Notes = make([]string, len(child.Notes))
					copy(note.Notes, child.Notes)
					if child.XMLName.Local == "sevenChrNote" {
						for _, def := range n.SevenChrDef {
							note.Notes = append(note.Notes, def.String())
						}
					}
					diag.Notes = append(diag.Notes, note)
				}
			}

			diags <- diag
			fallthrough
		default:
			for _, child := range n.Children {
				if len(child.Notes) != 0 {
					continue
				}

				child.walk(ctx, diags, depth+1, &n)
			}
		}
	}
}

var stopWords = map[string]bool{
	"code":       true,
	"to":         true,
	"identify":   true,
	"if":         true,
	"applicable": true,
	"for":        true,
	"of":         true,
	"the":        true,
	"or":         true,
	"associated": true,
	"any":        true,
	"and":        true,
	"as":         true,
	"such":       true,
	"specify":    true,
}

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	docDb, err := bolt.Open(boltFilename, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer docDb.Close()

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

	docDb.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(boltBucket))
		return nil
	})

	tx, err := docDb.Begin(true)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "docDb.Begin"))
	}
	defer tx.Commit()

	bucket, err := tx.CreateBucket([]byte("tabular"))
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "tx.CreateBucket"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	idx := 0
	for diag := range node.Walk(ctx) {
		docId := make([]byte, 8)
		docIdLen := binary.PutUvarint(docId, uint64(idx))

		doc, err := diag.MarshalMsg(nil)
		if err != nil {
			log.Fatalf("%+v\n", errors.Wrap(err, "diag.MarshalMsg"))
		}

		err = bucket.Put(docId[:docIdLen], doc)
		if err != nil {
			log.Fatalf("%+v\n", errors.Wrap(err, "bucket.Put"))
		}
		idx++
	}

	log.Println("parsed", idx, "diagnoses")
}
