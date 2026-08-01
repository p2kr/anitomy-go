package anitomy

func IsSpace(ch rune) bool {
	switch ch {
	case ' ', '\t', '\u00A0', '\u200B', '\u3000':
		return true
	default:
		return false
	}
}

func IsDash(ch rune) bool {
	switch ch {
	case '-', '\u00AD', '\u2010', '\u2011', '\u2012', '\u2013', '\u2014', '\u2015':
		return true
	default:
		return false
	}
}

func IsDelimiter(ch rune) bool {
	switch ch {
	case '_', '.', ',', '&', '~', '+', '|':
		return true
	default:
		return IsSpace(ch) || IsDash(ch)
	}
}
