package main

import (
	"log"
	"time"

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

type Index map[string]map[string]bool

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
	db, err := bolt.Open("../documents.db", 0600, &bolt.Options{Timeout: time.Second})
	fatalErr(err)
	defer db.Close()

	log.Println("parsing alphabetic terms...")
	n, err := ParseAlphabetic(db)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "ParseAlphabetic"))
	}
	log.Println("parsed", n, "alphabetic terms")

	log.Println("indexing alphabetic terms...")
	err = IndexAlphabetic(db)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "IndexAlphabetic"))
	}

	log.Println("parsing tabular diagnoses...")
	n, err = ParseTabular(db)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "ParseTabular"))
	}
	log.Println("parsed", n, "tabular diagnoses")

	log.Println("indexing tabular diagnoses...")
	err = IndexTabular(db)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "IndexTabular"))
	}
}
