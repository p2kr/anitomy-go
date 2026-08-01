package anitomy

func IsOpenBracket(ch rune) bool {
	switch ch {
	case '(', '[', '{', '\u300C', '\u300E', '\u3010', '\uFF08', '\uFF3B', '\uFF5B':
		return true
	default:
		return false
	}
}

func IsCloseBracket(ch rune) bool {
	switch ch {
	case ')', ']', '}', '\u300D', '\u300F', '\u3011', '\uFF09', '\uFF3D', '\uFF5D':
		return true
	default:
		return false
	}
}

func IsBracket(ch rune) bool {
	return IsOpenBracket(ch) || IsCloseBracket(ch)
}
