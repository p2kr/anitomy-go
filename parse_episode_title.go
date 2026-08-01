package anitomy

func ParseEpisodeTitle(tokens []Token) *Element {
	epIdx := -1
	for i := 0; i < len(tokens); i++ {
		if tokens[i].ElementKind != nil && *tokens[i].ElementKind == Episode {
			epIdx = i
			break
		}
	}
	if epIdx == -1 {
		return nil
	}

	first, last := -1, -1
	for i := epIdx + 1; i < len(tokens); i++ {
		if IsFreeToken(tokens[i]) && !IsEnclosedToken(tokens[i]) {
			first = i
			break
		}
	}

	isInvalidToken := func(t Token) bool {
		if !IsIdentifiedToken(t) {
			return false
		}
		if *t.ElementKind != Episode && *t.ElementKind != ReleaseVersion && *t.ElementKind != Season {
			return true
		}
		return false
	}

	if first != -1 {
		for i := epIdx; i < first; i++ {
			if isInvalidToken(tokens[i]) {
				first = -1
				break
			}
		}
	}

	if first != -1 {
		last = len(tokens)
		for i := first; i < len(tokens); i++ {
			t := tokens[i]
			isPartToken := IsIdentifiedToken(t) && *t.ElementKind == Part
			if IsOpenBracketToken(t) || (IsIdentifiedToken(t) && !isPartToken) {
				last = i
				break
			}
		}
	} else {
		// Fallback for corner brackets
		for i := epIdx + 1; i < len(tokens); i++ {
			if IsOpenBracketToken(tokens[i]) && tokens[i].Value == "「" {
				first = i + 1
				break
			}
		}
		if first != -1 {
			for i := first; i < len(tokens); i++ {
				if IsCloseBracketToken(tokens[i]) && tokens[i].Value == "」" {
					last = i
					break
				}
			}
			if last == -1 {
				return nil
			}
			for i := first; i < last; i++ {
				if IsIdentifiedToken(tokens[i]) {
					return nil
				}
			}
		}
	}

	if first == -1 || first >= last {
		return nil
	}

	span := tokens[first:last]
	val := BuildElementValue(span, KeepDelimitersNo)
	if val == "" || len(val) == 1 {
		return nil
	}

	for i := first; i < last; i++ {
		k := EpisodeTitle
		tokens[i].ElementKind = &k
	}

	return &Element{Kind: EpisodeTitle, Value: val, Position: span[0].Position}
}
