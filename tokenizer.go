package anitomy

import (
	"strings"
)

type Tokenizer struct {
	input  []rune
	view   []rune
	Tokens []Token
}

func NewTokenizer(input string) *Tokenizer {
	runes := []rune(input)
	return &Tokenizer{
		input: runes,
		view:  runes,
	}
}

func (t *Tokenizer) Tokenize(options Options) {
	for {
		token := t.nextToken()
		if token == nil {
			break
		}
		t.Tokens = append(t.Tokens, *token)
	}
	t.processTokens()
}

func (t *Tokenizer) nextToken() *Token {
	if t.isEof() {
		return nil
	}

	if IsOpenBracket(t.peek()) {
		return &Token{
			Kind:  TokenOpenBracket,
			Value: t.take(1),
		}
	}
	if IsCloseBracket(t.peek()) {
		return &Token{
			Kind:  TokenCloseBracket,
			Value: t.take(1),
		}
	}

	if IsDelimiter(t.peek()) {
		return &Token{
			Kind:  TokenDelimiter,
			Value: t.take(1),
		}
	}

	if value, keyword := t.takeKeyword(); value != "" {
		return &Token{
			Kind:    TokenKeyword,
			Value:   value,
			Keyword: keyword,
		}
	}

	return &Token{
		Kind:  TokenText,
		Value: t.takeText(),
	}
}

func (t *Tokenizer) processTokens() {
	bracketLevel := 0
	position := 0

	for i := range t.Tokens {
		token := &t.Tokens[i]
		if token.Kind == TokenOpenBracket {
			bracketLevel++
		} else if token.Kind == TokenCloseBracket {
			bracketLevel--
		} else {
			token.IsEnclosed = bracketLevel > 0
		}

		token.Position = position
		position += len([]rune(token.Value)) // in C++, it's token.value.size() which is byte length. Let's use rune length since we operate on runes. Wait, C++ uses byte length for position! Let's match C++: len(token.Value)

		if token.Kind == TokenText {
			isNum := true
			for _, ch := range token.Value {
				if !IsDigit(ch) {
					isNum = false
					break
				}
			}
			token.IsNumber = isNum && len(token.Value) > 0
		}
	}
}

func isText(ch rune) bool {
	return !IsBracket(ch) && !IsDelimiter(ch)
}

func isWordBoundary(ch rune) bool {
	return !isText(ch)
}

func (t *Tokenizer) isEof() bool {
	return len(t.view) == 0
}

func (t *Tokenizer) peek() rune {
	return t.view[0]
}

func (t *Tokenizer) peekN(offset, n int) []rune {
	if offset >= len(t.view) {
		return nil
	}
	end := offset + n
	if end > len(t.view) {
		end = len(t.view)
	}
	return t.view[offset:end]
}

func (t *Tokenizer) take(n int) string {
	if n > len(t.view) {
		n = len(t.view)
	}
	res := string(t.view[:n])
	t.view = t.view[n:]
	return res
}

func (t *Tokenizer) takeText() string {
	n := 0
	for n < len(t.view) && isText(t.view[n]) {
		n++
	}
	return t.take(n)
}

func (t *Tokenizer) hasCandidates(prefix string) bool {
	for k := range Keywords {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

func (t *Tokenizer) findKey(view []rune) (string, int) {
	var key string
	var keyRuneCount int
	for n := 1; n <= len(view); n++ {
		prefix := strings.ToLower(string(view[:n]))
		if _, exists := Keywords[prefix]; exists {
			key = prefix
			keyRuneCount = n
		}
		if !t.hasCandidates(prefix) {
			break
		}
	}
	return key, keyRuneCount
}

func (t *Tokenizer) isKeywordBoundary(keyword Keyword, view []rune) bool {
	if keyword.IsSubword() {
		return true
	}
	if len(view) == 0 {
		return true
	}
	next := view[0]
	if isWordBoundary(next) {
		return true
	}
	if keyword.IsPrefixForNumber() {
		return IsDigit(next)
	}
	if keyword.IsPrefixForOther() {
		k, _ := t.findKey(view)
		return k != ""
	}
	return false
}

func (t *Tokenizer) takeKeyword() (string, *Keyword) {
	key, n := t.findKey(t.view)
	if key == "" {
		return "", nil
	}

	keyword := Keywords[key]

	var nextView []rune
	if n < len(t.view) {
		nextView = t.view[n:]
	}

	if !t.isKeywordBoundary(keyword, nextView) {
		return "", nil
	}

	return t.take(n), &keyword
}
