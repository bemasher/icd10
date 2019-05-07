package util

import (
	"strings"
	"unicode"

	"github.com/dchest/stemmer/porter2"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func normStr(src string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	dst, _, _ := transform.Bytes(t, []byte(src))

	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}

		switch r {
		case '.', '-', '*':
			return r
		case '\'':
			return -1
		}

		return ' '
	}, string(dst))
}

type Token struct {
	Src        string
	Forms      []string
	IsPrefix   bool
	IsNegative bool
}

func NewToken(src string) (t Token) {
	t.Src = src
	t.IsPrefix = strings.HasSuffix(src, "*")
	t.IsNegative = strings.HasPrefix(src, "-")
	src = strings.TrimSuffix(src, "*")
	src = strings.TrimPrefix(src, "-")

	forms := map[string]bool{}

	t.AddForm(forms, src)
	t.AddForm(forms, porter2.Stemmer.Stem(src))

	src = strings.Map(func(r rune) rune {
		if r == '-' {
			return -1
		}
		return r
	}, src)

	t.AddForm(forms, src)
	t.AddForm(forms, porter2.Stemmer.Stem(src))

	return
}

func (t *Token) AddForm(formSet map[string]bool, form string) {
	if form == "" || formSet[form] {
		return
	}
	t.Forms = append(t.Forms, form)
	formSet[form] = true
}

func Tokenize(s string) (tokens []Token) {
	fields := strings.Fields(normStr(s))
	tokens = make([]Token, len(fields))
	for idx, token := range fields {
		tokens[idx] = NewToken(token)
	}
	return
}
