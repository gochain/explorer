package utils

import (
	"strings"
	"unicode/utf8"
)

func fixUtf(r rune) rune {
	if r == utf8.RuneError || r == '\u0000' {
		return -1
	}
	return r
}

//CleanUpText removes from the text all non-utf and null characters
func CleanUpText(t string) string {
	return strings.Map(fixUtf, t)
}
