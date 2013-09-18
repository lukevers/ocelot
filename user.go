package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type User struct {
	Address     string
	Name        string
	Description string
	Website     string
}

// CreateFullName is a func that takes two strings and returns a full
// name from the first name and last name. The first letter of each is
// capitalized and the rest are lowered.
func CreateFullName(f, l string) (name string) {
	f = strings.ToLower(strings.Trim(f, " "))
	l = strings.ToLower(strings.Trim(l, " "))
	
	name = UpperFirst(f) + " " + UpperFirst(l)
	
	return
}

// UpperFirst takes a string and returns the same string with the
// first letter of the string uppercased.
func UpperFirst(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}