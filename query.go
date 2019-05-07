package main

import (
	"sort"
	"strings"

	"github.com/bemasher/icd10/util"
	"github.com/pkg/errors"
)

type QueryFn func(string) (interface{}, error)

var querySets = map[string][]Query{
	"cm": {
		Query{"Alphabetic", QueryAlphabetic},
		Query{"Tabular", QueryTabular},
	},
}

type Query struct {
	Name string
	Fn   QueryFn
}

func QueryAlphabetic(qry string) (results interface{}, err error) {
	var terms []util.AlphabeticTerm

	err = SearchDocs("alphabetic", "alphabetic_index", qry, func(doc []byte) error {
		var term util.AlphabeticTerm
		_, err = term.UnmarshalMsg(doc)
		if err != nil {
			return errors.Wrap(err, "term.UnmarshalMsg")
		}
		terms = append(terms, term)
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "SearchDocs")
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

func QueryTabular(qry string) (results interface{}, err error) {
	var diags []util.Diag

	err = SearchDocs("tabular", "tabular_index", qry, func(doc []byte) error {
		var diag util.Diag
		_, err = diag.UnmarshalMsg(doc)
		if err != nil {
			return errors.Wrap(err, "term.UnmarshalMsg")
		}
		diags = append(diags, diag)
		return nil
	})

	sort.Slice(diags, func(i, j int) bool {
		return strings.Compare(diags[i].Code, diags[j].Code) < 0
	})

	trim := ResultLimit
	if len(diags) < trim {
		trim = len(diags)
	}

	return diags[:trim], nil
}
