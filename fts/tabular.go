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

type TabularNode struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`

	Code string `xml:"name"`
	Desc string `xml:"desc"`

	Notes       []string    `xml:"note"`
	SevenChrDef []Extension `xml:"sevenChrDef>extension"`

	Children []TabularNode `xml:",any"`
}

type Extension struct {
	Char string `xml:"char,attr"`
	Def  string `xml:",innerxml"`
}

func (e Extension) String() string {
	return fmt.Sprintf("%q - %s", e.Char, e.Def)
}

func (n TabularNode) Walk(ctx context.Context) chan util.Diag {
	diags := make(chan util.Diag)

	go func() {
		for _, child := range n.Children {
			child.walk(ctx, diags, 0, nil)
		}
		close(diags)
	}()

	return diags
}

func (n TabularNode) walk(ctx context.Context, diags chan util.Diag, depth int, parent *TabularNode) {
	if len(n.Children) == 0 {
		n.Children = nil
	}

	select {
	case <-ctx.Done():
		return
	default:
		if n.XMLName.Local == "diag" {
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
		}

		for _, child := range n.Children {
			if len(child.Notes) != 0 {
				continue
			}

			child.walk(ctx, diags, depth+1, &n)
		}
	}
}

const (
	tabularBkt    = "tabular"
	tabularIdxBkt = "tabular_index"
)

func ParseTabular(db *bolt.DB) (n uint64, err error) {
	xmlFile, err := os.Open("icd10cm_tabular_2019.xml")
	if err != nil {
		return 0, errors.Wrap(err, "os.Open")
	}
	defer xmlFile.Close()

	xmlDecoder := xml.NewDecoder(xmlFile)

	var node TabularNode
	err = xmlDecoder.Decode(&node)
	if err != nil {
		return 0, errors.Wrap(err, "xmlDecoder.Decode")
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(tabularBkt))
		return nil
	})

	tx, err := db.Begin(true)
	if err != nil {
		return 0, errors.Wrap(err, "docDb.Begin")
	}
	defer tx.Commit()

	bucket, err := tx.CreateBucket([]byte(tabularBkt))
	if err != nil {
		return 0, errors.Wrap(err, "tx.CreateBucket")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	idx := uint64(0)
	for diag := range node.Walk(ctx) {
		docId := make([]byte, 8)
		docIdLen := binary.PutUvarint(docId, uint64(idx))

		doc, err := diag.MarshalMsg(nil)
		if err != nil {
			return 0, errors.Wrap(err, "diag.MarshalMsg")
		}

		err = bucket.Put(docId[:docIdLen], doc)
		if err != nil {
			return 0, errors.Wrap(err, "bucket.Put")
		}
		idx++
	}

	return idx, nil
}

func IndexTabular(db *bolt.DB) (err error) {
	index := Index{}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tabularBkt))
		if bucket == nil {
			return errors.New("tabular bucket missing")
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var diag util.Diag
			_, err := diag.UnmarshalMsg(v)
			if err != nil {
				log.Printf("%q\n", string(v))
				return errors.Wrap(err, "diag.UnmarshalMsg")
			}

			indexPhrase(index, util.Tokenize, string(k), diag.Code)
			indexPhrase(index, util.Tokenize, string(k), diag.Desc)

			for _, noteGroup := range diag.Notes {
				switch noteGroup.Kind {
				case "excludes1", "excludes2":
					continue
				}
				for idx, note := range noteGroup.Notes {
					if noteGroup.Kind == "sevenChrNote" && idx == 0 {
						continue
					}
					indexPhrase(index, util.Tokenize, string(k), note)
				}
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "db.View")
	}

	err = WriteIndex(db, "tabular_index", index)
	return errors.Wrap(err, "WriteIndex")
}
