package anitomy

import (
	"strconv"
	"strings"
)

var ordinalTable = map[string]string{
	"1st": "1", "First": "1",
	"2nd": "2", "Second": "2",
	"3rd": "3", "Third": "3",
	"4th": "4", "Fourth": "4",
	"5th": "5", "Fifth": "5",
	"6th": "6", "Sixth": "6",
	"7th": "7", "Seventh": "7",
	"8th": "8", "Eighth": "8",
	"9th": "9", "Ninth": "9",
}

func FromOrdinalNumber(input string) string {
	if val, ok := ordinalTable[input]; ok {
		return val
	}
	return ""
}

var romanTable = map[string]string{
	"II":  "2",
	"III": "3",
	"IV":  "4",
}

func FromRomanNumber(input string) string {
	if val, ok := romanTable[input]; ok {
		return val
	}
	return ""
}

func IsAlpha(ch rune) bool {
	return ('A' <= ch && ch <= 'Z') || ('a' <= ch && ch <= 'z')
}

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func IsXDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('A' <= ch && ch <= 'F') || ('a' <= ch && ch <= 'f')
}

func ToInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

func ToFloat(str string) float32 {
	val, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0.0
	}
	return float32(val)
}

func ToLower(ch rune) rune {
	if 'A' <= ch && ch <= 'Z' {
		return ch + ('a' - 'A')
	}
	return ch
}

func EqualTo(a, b byte) bool {
	return ToLower(rune(a)) == ToLower(rune(b))
}

func Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}

func FindAllIf[T any](slice []T, predicate func(T) bool) []T {
	var found []T
	for _, item := range slice {
		if predicate(item) {
			found = append(found, item)
		}
	}
	return found
}
