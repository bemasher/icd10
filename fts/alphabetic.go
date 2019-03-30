package main

import (
	"context"
	"encoding/binary"
	"encoding/xml"
	"os"
	"strings"

	"github.com/bemasher/icd10/util"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type AlphabeticNode struct {
	XMLName  xml.Name
	Content  string           `xml:",chardata"`
	Children []AlphabeticNode `xml:",any"`

	Title util.Title `xml:"title"`
}

func (n AlphabeticNode) Walk(ctx context.Context) chan util.AlphabeticTerm {
	terms := make(chan util.AlphabeticTerm)

	go func() {
		defer close(terms)
		for _, child := range n.Children {
			select {
			case <-ctx.Done():
				return
			default:
				child.walk(ctx, terms, nil)
			}
		}
	}()

	return terms
}

func (n AlphabeticNode) walk(ctx context.Context, terms chan util.AlphabeticTerm, breadcrumbs []string) {
	if len(n.Children) == 0 {
		n.Children = nil
	}

	select {
	case <-ctx.Done():
		return
	default:
		switch n.XMLName.Local {
		case "mainTerm", "term":
			if len(n.Title.Nemod) > 0 {
				n.Title.Title += " " + n.Title.Nemod
			}
			breadcrumbs = append(breadcrumbs, n.Title.Title)

			var term util.AlphabeticTerm

			term.Title = strings.Join(breadcrumbs, ", ")

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
		}

		for _, child := range n.Children {
			child.walk(ctx, terms, breadcrumbs)
		}
	}
}

const (
	alphaBkt    = "alphabetic"
	alphaIdxBkt = "alphabetic_index"
)

func ParseAlphabetic(db *bolt.DB) (n int, err error) {
	xmlFile, err := os.Open("icd10cm_index_2019.xml")
	if err != nil {
		return 0, errors.Wrap(err, "os.Open")
	}
	defer xmlFile.Close()

	xmlDecoder := xml.NewDecoder(xmlFile)

	var node AlphabeticNode
	err = xmlDecoder.Decode(&node)
	if err != nil {
		return 0, errors.Wrap(err, "xmlDecoder.Decode")
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(alphaBkt))
		return nil
	})

	tx, err := db.Begin(true)
	if err != nil {
		return 0, errors.Wrap(err, "db.Begin")
	}
	defer tx.Commit()

	bucket, err := tx.CreateBucket([]byte(alphaBkt))
	if err != nil {
		return 0, errors.Wrap(err, "tx.CreateBucket")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	idx := 0
	for term := range node.Walk(ctx) {
		docId := make([]byte, 8)
		docIdLen := binary.PutUvarint(docId, uint64(idx))

		doc, err := term.MarshalMsg(nil)
		if err != nil {
			return 0, errors.Wrap(err, "term.MarshalMsg")
		}

		bucket.Put(docId[:docIdLen], doc)

		idx++
	}

	return idx, nil
}

func IndexAlphabetic(db *bolt.DB) (err error) {
	index := map[string]map[string]bool{}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(alphaBkt))
		if bucket == nil {
			return errors.New("alphabetic bucket missing")
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var term util.AlphabeticTerm
			_, err := term.UnmarshalMsg(v)
			if err != nil {
				fatalErr(errors.Wrap(err, "term.UnmarshalMsg"))
			}

			indexPhrase(index, util.Tokenize, string(k), term.Code)
			indexPhrase(index, util.Tokenize, string(k), term.Manif)
			indexPhrase(index, util.Tokenize, string(k), term.Title)

			for _, attr := range term.Attrs {
				indexPhrase(index, util.Tokenize, string(k), attr.Value)
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "db.View")
	}

	db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(alphaIdxBkt))
	})

	err = WriteIndex(db, alphaIdxBkt, index)
	return errors.Wrap(err, "WriteIndex")
}
