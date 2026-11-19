package main

import (
	"strings"
	"unicode/utf8"
)

func structToFilename(s string) string {
	if s == "" {
		return s
	}

	_, size := utf8.DecodeRuneInString(s)
	return strings.ToLower(s[:size]) + s[size:] + ".go"
}
