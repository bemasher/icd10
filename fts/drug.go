package main

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/bemasher/icd10/util"
)

type DrugNode struct {
	XMLName  xml.Name
	Content  string     `xml:",chardata"`
	Attrs    []xml.Attr `xml:",attr,any"`
	Children []DrugNode `xml:",any"`

	Title   util.Title `xml:"title"`
	Codes   []string   `xml:"cell"`
	See     string     `xml:"see"`
	SeeAlso string     `xml:"seeAlso"`
}

func (n DrugNode) Walk(ctx context.Context) chan util.AlphabeticTerm {
	terms := make(chan util.AlphabeticTerm)

	go func() {
		for _, child := range n.Children {
			child.walk(ctx, terms, 0, nil)
		}
		close(terms)
	}()

	return terms
}

func (n DrugNode) walk(ctx context.Context, terms chan util.AlphabeticTerm, depth int, breadcrumbs []string) {
	select {
	case <-ctx.Done():
		return
	default:
		switch n.XMLName.Local {
		case "indexHeading":

		case "mainTerm", "term":
			if len(n.Title.Nemod) > 0 {
				n.Title.Title += " " + n.Title.Nemod
			}

			breadcrumbs = append(breadcrumbs, n.Title.Title)
			title := strings.Join(breadcrumbs, ", ")

			for idx, code := range n.Codes {
				if code == "--" {
					continue
				}

				var term util.AlphabeticTerm

				term.Code = code

				term.Title = title + ", " + drugHeaders[idx]

				if n.See != "" {
					term.Attrs = append(term.Attrs, util.Attr{"see", n.See})
				}
				if n.SeeAlso != "" {
					term.Attrs = append(term.Attrs, util.Attr{"seeAlso", n.SeeAlso})
				}

				terms <- term
			}
		}

		for _, child := range n.Children {
			child.walk(ctx, terms, depth+1, breadcrumbs)
		}
	}
}

var drugHeaders = []string{
	"Poisoning Accidental (unintentional)",
	"Poisoning Intentional self-harm",
	"Poisoning Assault",
	"Poisoning Undetermined",
	"Adverse effect",
	"Underdosing",
}
