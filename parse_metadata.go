package anitomy

func ParseFileExtension(tokens []Token) *Element {
	// Simple implementation
	extensions := map[string]bool{
		"3gp": true, "avi": true, "divx": true, "flv": true, "m2ts": true,
		"m4v": true, "mkv": true, "mov": true, "mp4": true, "mpg": true,
		"ogm": true, "rm": true, "rmvb": true, "ts": true, "webm": true, "wmv": true,
	}

	if len(tokens) < 2 {
		return nil
	}

	lastToken := &tokens[len(tokens)-1]
	prevToken := &tokens[len(tokens)-2]

	isExt := (IsKeywordToken(*lastToken) || IsTextToken(*lastToken)) && extensions[lastToken.Value]
	isDot := IsDelimiterToken(*prevToken) && prevToken.Value == "."

	if !isExt || !isDot {
		return nil
	}

	lastToken.Kind = TokenText
	lastToken.Keyword = nil
	k := FileExtension
	lastToken.ElementKind = &k

	e := ElementFromToken(FileExtension, *lastToken, "", -1)
	return &e
}

func ParseFileChecksum(tokens []Token) *Element {
	for i := len(tokens) - 1; i >= 0; i-- {
		t := &tokens[i]
		if IsFreeToken(*t) && len(t.Value) == 8 {
			allXDigit := true
			for _, ch := range t.Value {
				if !IsXDigit(ch) {
					allXDigit = false
					break
				}
			}
			if allXDigit {
				k := FileChecksum
				t.ElementKind = &k
				e := ElementFromToken(FileChecksum, *t, "", -1)
				return &e
			}
		}
	}
	return nil
}

func ParseReleaseGroup(tokens []Token) *Element {
	// e.g. [Group] Title ...
	var first, last int = -1, -1
	for i := 0; i < len(tokens); i++ {
		if IsEnclosedToken(tokens[i]) && !IsIdentifiedToken(tokens[i]) {
			first = i
			break
		}
	}
	if first != -1 {
		for i := first; i < len(tokens); i++ {
			if IsCloseBracketToken(tokens[i]) || IsIdentifiedToken(tokens[i]) {
				last = i
				break
			}
		}
	}

	if first != -1 && last != -1 {
		prev := FindPrevToken(tokens, first, IsNotDelimiterToken)
		if prev != -1 && !IsOpenBracketToken(tokens[prev]) {
			first, last = -1, -1 // skip
		} else if !IsCloseBracketToken(tokens[last]) {
			first, last = -1, -1 // skip
		}
	}

	if first == -1 {
		// fall back to last token before file extension
		idx := -1
		for i := len(tokens) - 1; i >= 0; i-- {
			if tokens[i].ElementKind == nil || *tokens[i].ElementKind != FileExtension {
				if IsNotDelimiterToken(tokens[i]) {
					idx = i
					break
				}
			}
		}
		if idx != -1 && IsFreeToken(tokens[idx]) {
			if idx-1 >= 0 && IsDelimiterToken(tokens[idx-1]) && IsDashToken(tokens[idx-1]) {
				first = idx
				last = idx + 1
			}
		}
	}

	if first == -1 || last == -1 || first >= last {
		return nil
	}

	span := tokens[first:last]
	value := BuildElementValue(span, KeepDelimitersYes)
	if value == "" {
		return nil
	}

	for i := first; i < last; i++ {
		k := ReleaseGroup
		tokens[i].ElementKind = &k
	}

	return &Element{Kind: ReleaseGroup, Value: value, Position: span[0].Position}
}
