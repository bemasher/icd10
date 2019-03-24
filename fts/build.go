package main

import (
	"log"

	"github.com/bemasher/icd10/util"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

func fatalErr(err error) {
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func indexPhrase(index map[string]map[string]bool, tokenFn func(string) []string, docId, phrase string) {
	for _, token := range tokenFn(phrase) {
		if _, ok := index[token]; !ok {
			index[token] = map[string]bool{}
		}
		index[token][docId] = true
	}
}

func buildAlphabeticIndex(db *bolt.DB) {
	index := map[string]map[string]bool{}

	log.Println("Building alphabetic index...")
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("alphabetic"))
		if bucket == nil {
			return errors.New("alphabetic bucket missing")
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var term util.Term
			_, err := term.UnmarshalMsg(v)
			if err != nil {
				fatalErr(errors.Wrap(err, "term.UnmarshalMsg"))
			}

			indexPhrase(index, util.TokenizeCode, string(k), term.Code)
			indexPhrase(index, util.TokenizeCode, string(k), term.Manif)
			indexPhrase(index, util.Tokenize, string(k), term.Title)
			indexPhrase(index, util.Tokenize, string(k), term.Nemod)

			for _, attr := range term.Attrs {
				indexPhrase(index, util.Tokenize, string(k), attr.Value)
			}
		}

		return nil
	})
	fatalErr(errors.Wrap(err, "db.View"))

	db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte("alphabetic_index"))
	})

	log.Println("Writing alphabetic index...")
	err = WriteIndex(db, "alphabetic_index", index)
	fatalErr(errors.Wrap(err, "WriteIndex"))
}

func buildTabularIndex(db *bolt.DB) {
	index := map[string]map[string]bool{}

	log.Println("Building tabular index...")
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tabular"))
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

			indexPhrase(index, util.TokenizeCode, string(k), diag.Code)
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
	fatalErr(errors.Wrap(err, "db.View"))

	log.Println("Writing tabular index...")
	err = WriteIndex(db, "tabular_index", index)
	fatalErr(errors.Wrap(err, "WriteIndex"))
}

func WriteIndex(db *bolt.DB, bucketName string, index map[string]map[string]bool) (err error) {
	db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return errors.Wrap(err, "tx.CreateBucket")
		}

		for term, codes := range index {
			doc, err := util.DocIDMap(codes).MarshalMsg(nil)
			if err != nil {
				return errors.Wrap(err, "DocIDMap.MarshalMsg")
			}

			bucket.Put([]byte(term), doc)
		}

		return nil
	})
	return errors.Wrap(err, "db.Update")
}

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	db, err := bolt.Open("../documents.db", 0600, nil)
	fatalErr(err)
	defer db.Close()

	buildAlphabeticIndex(db)
	buildTabularIndex(db)
}
