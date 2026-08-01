package anitomy

func ParseTitle(tokens []Token) *Element {
	if len(tokens) > 0 && tokens[0].ElementKind != nil && *tokens[0].ElementKind == Episode {
		return nil
	}

	first, last := -1, -1

	// Find the first free unenclosed range
	for i := 0; i < len(tokens); i++ {
		if IsFreeToken(tokens[i]) && !IsEnclosedToken(tokens[i]) {
			first = i
			break
		}
	}
	if first != -1 {
		last = first
		for i := first; i < len(tokens); i++ {
			if IsIdentifiedToken(tokens[i]) {
				break
			}
			last = i + 1
		}
	} else {
		// Fallback to second enclosed range
		for i := 0; i < len(tokens); i++ {
			if IsCloseBracketToken(tokens[i]) {
				for j := i + 1; j < len(tokens); j++ {
					if IsFreeToken(tokens[j]) {
						first = j
						break
					}
				}
				break
			}
		}
		if first != -1 {
			last = first
			for i := first; i < len(tokens); i++ {
				if IsBracketToken(tokens[i]) {
					break
				}
				last = i + 1
			}
		}
	}

	if first == -1 || first == last {
		return nil
	}

	// Mismatched brackets check
	var openBrackets []int
	var closeBrackets []int
	for i := first; i < last; i++ {
		if IsOpenBracketToken(tokens[i]) {
			openBrackets = append(openBrackets, i)
		} else if IsCloseBracketToken(tokens[i]) {
			closeBrackets = append(closeBrackets, i)
		}
	}
	if len(openBrackets) > 0 && len(openBrackets) != len(closeBrackets) {
		last = openBrackets[len(openBrackets)-1]
	}

	// Trailing brackets check
	lastIdx := last - 1
	if lastIdx >= 0 {
		prev := FindPrevToken(tokens, last, IsNotDelimiterToken)
		if prev >= 0 && IsCloseBracketToken(tokens[prev]) && tokens[prev].Value != ")" {
			openPrev := FindPrevToken(tokens, prev+1, IsOpenBracketToken)
			if openPrev >= 0 {
				last = openPrev
			}
		}
	}

	if first >= last {
		return nil
	}

	span := tokens[first:last]
	value := BuildElementValue(span, KeepDelimitersNo)
	if value == "" {
		return nil
	}

	for i := first; i < last; i++ {
		k := Title
		tokens[i].ElementKind = &k
	}

	return &Element{Kind: Title, Value: value, Position: span[0].Position}
}
