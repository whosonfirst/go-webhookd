package ucd

import (
	"fmt"
	unicodedata "github.com/cooperhewitt/go-ucd/unicodedata"
	unihan "github.com/cooperhewitt/go-ucd/unihan"
	"strings"
	"unicode/utf8"
)

type UCDName struct {
	Char string
	Hex  string
	Name string
}

func (u UCDName) String() string {
	return u.Name
}

func Name(char string) (f UCDName) {

	rune, _ := utf8.DecodeRuneInString(char)
	hex := fmt.Sprintf("%04X", rune)

	name, ok := unicodedata.UCD[hex]

	if ok == false {
		name, ok = unihan.UCDHan[hex]
	}

	if ok == false {
		hex = fmt.Sprintf("%05X", rune)
		name, ok = unihan.UCDHan[hex]
	}

	return UCDName{char, hex, name}
}

func NamesForString(s string) (n []UCDName) {

	chars := strings.Split(s, "")
	count := len(chars)

	results := make([]UCDName, count)

	for idx, char := range chars {
		name := Name(char)
		results[idx] = name
	}

	return results
}
