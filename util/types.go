package util

//go:generate msgp

//msgp:ignore Title
//msgp:tuple Term Attr Diag Note DrugTerm

// Alphabetic Index
type AlphabeticTerm struct {
	Title string
	Code  string
	Manif string
	Src   string

	Attrs []Attr
}

type Attr struct {
	Attr  string
	Value string
}

// Tabular Index
type Diag struct {
	Code string
	Desc string

	Notes []Note
}

type Note struct {
	Kind  string
	Notes []string
}

// FTS Types
type DocIDMap map[string]bool

// Utility Types
type Title struct {
	Title string `xml:",chardata"`
	Nemod string `xml:"nemod"`
}
