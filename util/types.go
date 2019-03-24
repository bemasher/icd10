package util

//go:generate msgp

// Alphabetic Index
type Term struct {
	Title string `msg:"t"`
	Nemod string `msg:"n"`
	Code  string `msg:"c"`
	Manif string `msg:"m"`

	Attrs []Attr `msg:"as"`
}

type Attr struct {
	Attr  string `msg:"a"`
	Value string `msg:"v"`
}

// Tabular Index
type Diag struct {
	Code string `msg:"c"`
	Desc string `msg:"d"`

	Notes []Note `msg:"n"`
}

type Note struct {
	Kind  string   `msg:"k"`
	Notes []string `msg:"ns"`
}

type DocIDMap map[string]bool
