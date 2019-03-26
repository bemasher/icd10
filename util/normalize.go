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

func Tokenize(s string) (tokens []string) {
	tokens = strings.Fields(normStr(s))
	for idx, token := range tokens {
		isPrefix := strings.HasSuffix(token, "*")
		token = strings.TrimRight(token, ".-*")

		stemmed := porter2.Stemmer.Stem(token)
		if isPrefix {
			stemmed += "*"
		}
		tokens[idx] = stemmed
	}
	return
}
