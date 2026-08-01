package anitomy

import (
	"strings"
)

type KeepDelimiters bool

const (
	KeepDelimitersNo  KeepDelimiters = false
	KeepDelimitersYes KeepDelimiters = true
)

func ElementFromToken(kind ElementKind, token Token, value string, position int) Element {
	val := value
	if val == "" {
		val = token.Value
	}
	pos := position
	if pos == -1 {
		pos = token.Position
	}
	return Element{
		Kind:     kind,
		Value:    val,
		Position: pos,
	}
}

func BuildElementValue(tokens []Token, keep KeepDelimiters) string {
	if len(tokens) == 0 {
		return ""
	}

	delimiters := make(map[rune]bool)
	for _, t := range tokens {
		if IsDelimiterToken(t) && len(t.Value) > 0 {
			firstRune := []rune(t.Value)[0]
			delimiters[firstRune] = true
		}
	}

	hasSingleDelimiter := len(delimiters) == 1
	hasSpaces := false
	for d := range delimiters {
		if IsSpace(d) {
			hasSpaces = true
			break
		}
	}
	hasUnderscores := delimiters['_']

	isTransformableDelimiter := func(token Token) bool {
		if keep == KeepDelimitersYes {
			return false
		}
		if IsNotDelimiterToken(token) {
			return false
		}
		if len(token.Value) == 0 {
			return false
		}
		ch := []rune(token.Value)[0]

		if ch == ',' || ch == '&' || ch == '~' {
			return false
		}
		if IsSpace(ch) || ch == '_' {
			return true
		}
		if hasSpaces || hasUnderscores {
			return false
		}
		if ch == '.' {
			return true
		}
		return hasSingleDelimiter
	}

	if keep == KeepDelimitersNo {
		var prevDelimiter rune
		for len(tokens) > 0 && IsDelimiterToken(tokens[len(tokens)-1]) {
			t := tokens[len(tokens)-1]
			if len(t.Value) > 0 {
				delimiter := []rune(t.Value)[0]
				if delimiter == '~' {
					break
				}
				if delimiter == '.' && IsSpace(prevDelimiter) {
					break
				}
				prevDelimiter = delimiter
			}
			tokens = tokens[:len(tokens)-1]
		}
	}

	var elementValue strings.Builder
	for _, token := range tokens {
		if isTransformableDelimiter(token) {
			elementValue.WriteRune(' ')
		} else {
			elementValue.WriteString(token.Value)
		}
	}

	return elementValue.String()
}
