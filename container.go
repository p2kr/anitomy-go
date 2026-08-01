package anitomy

func FindPrevToken(tokens []Token, index int, predicate func(Token) bool) int {
	if index <= 0 {
		return -1
	}
	for i := index - 1; i >= 0; i-- {
		if predicate(tokens[i]) {
			return i
		}
	}
	return -1
}

func FindNextToken(tokens []Token, index int, predicate func(Token) bool) int {
	if index < 0 || index >= len(tokens)-1 {
		return -1
	}
	for i := index + 1; i < len(tokens); i++ {
		if predicate(tokens[i]) {
			return i
		}
	}
	return -1
}
