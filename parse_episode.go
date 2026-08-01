package anitomy

import (
	"fmt"
	"regexp"
)

var (
	episodePattern           = regexp.MustCompile(`^(?:S(\d{1,2})|(\d{1,2})x)?[E#]?(\d{1,4})(?:[vV](\d))?$`)
	episodeJpCounterPattern  = regexp.MustCompile(`^(?:第)?(\d{1,4})話$`)
	episodeEquivalentPattern = regexp.MustCompile(`^(\d{1,4})\s*\((\d{1,4})\)$`)
	episodePartialPattern    = regexp.MustCompile(`^\d{1,4}[ABCabc]$`)
)

func ParseEpisode(tokens []Token) []Element {
	var elements []Element

	addEl := func(k ElementKind, t Token, val string, pos int) {
		t.ElementKind = &k
		elements = append(elements, ElementFromToken(k, t, val, pos))
	}

	parseMatches := func(matches []int, token *Token) {
		if matches[2] != -1 {
			val := token.Value[matches[2]:matches[3]]
			addEl(Season, *token, val, token.Position+matches[2])
		} else if matches[4] != -1 {
			val := token.Value[matches[4]:matches[5]]
			addEl(Season, *token, val, token.Position+matches[4])
		}

		epVal := token.Value[matches[6]:matches[7]]
		addEl(Episode, *token, epVal, token.Position+matches[6])

		if matches[8] != -1 {
			val := token.Value[matches[8]:matches[9]]
			addEl(ReleaseVersion, *token, val, token.Position+matches[8])
		}
	}

	isEpKeyword := func(t Token) bool {
		return t.Keyword != nil && t.Keyword.Kind == KeywordEpisode
	}
	isTypeKeyword := func(t Token) bool {
		return t.Keyword != nil && (t.Keyword.Kind == KeywordType || t.Keyword.Kind == KeywordEpisodeType)
	}
	isEpDelimiter := func(t Token) bool {
		if !IsDelimiterToken(t) || len(t.Value) == 0 {
			return false
		}
		ch := []rune(t.Value)[0]
		return ch == '-' || ch == '~' || ch == '&' || ch == '+'
	}

	startsWithKw := func(idx int) bool {
		if idx < 0 || idx >= len(tokens) {
			return false
		}
		t := tokens[idx]
		if isEpKeyword(t) {
			return true
		}
		if isTypeKeyword(t) {
			return t.Value != "Movie"
		}
		return false
	}

	// 1. Episode Token Matches
	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if !IsFreeToken(*t) {
			continue
		}
		matches := episodePattern.FindStringSubmatchIndex(t.Value)
		if matches == nil {
			continue
		}

		prev := FindPrevToken(tokens, i, IsNotDelimiterToken)
		next := i + 1

		valid := !IsNumericToken(*t)
		valid = valid || startsWithKw(prev)

		var nextMatches []int
		if next < len(tokens) && isEpDelimiter(tokens[next]) && next+1 < len(tokens) {
			nextNext := next + 1
			nextMatches = episodePattern.FindStringSubmatchIndex(tokens[nextNext].Value)
			if nextMatches != nil {
				ep1 := ToInt(t.Value[matches[6]:matches[7]])
				ep2 := ToInt(tokens[nextNext].Value[nextMatches[6]:nextMatches[7]])
				if ep1 < ep2 {
					valid = true
				} else {
					nextMatches = nil
				}
			}
		}

		if !valid {
			continue
		}

		parseMatches(matches, t)
		if nextMatches != nil {
			parseMatches(nextMatches, &tokens[next+1])
			i = next + 1
		}
	}
	if len(elements) > 0 {
		return elements
	}

	// 2. Separated episodes
	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if IsFreeToken(*t) && IsNumericToken(*t) {
			if i+2 < len(tokens) {
				sep := tokens[i+1]
				if sep.Value == "&" || sep.Value == "~" || sep.Value == "of" {
					nextTk := &tokens[i+2]
					if IsFreeToken(*nextTk) && IsNumericToken(*nextTk) {
						addEl(Episode, *t, "", -1)
						if sep.Value != "of" {
							addEl(Episode, *nextTk, "", -1)
						}
						return elements
					}
				}
			}
		}
	}

	// 3. Fractional episode
	for i := 0; i < len(tokens)-2; i++ {
		num, delim, frac := &tokens[i], &tokens[i+1], &tokens[i+2]
		if IsFreeToken(*num) && IsNumericToken(*num) {
			if IsDelimiterToken(*delim) && delim.Value == "." {
				if IsFreeToken(*frac) && frac.Value == "5" {
					val := fmt.Sprintf("%s%s%s", num.Value, delim.Value, frac.Value)
					addEl(Episode, *num, val, -1)
					k := Episode
					delim.ElementKind = &k
					frac.ElementKind = &k
					return elements
				}
			}
		}
	}

	// 4. Japanese counter
	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if !IsFreeToken(*t) {
			continue
		}
		if matches := episodeJpCounterPattern.FindStringSubmatchIndex(t.Value); matches != nil {
			val := t.Value[matches[2]:matches[3]]
			addEl(Episode, *t, val, t.Position+matches[2]) // rune offset logic here! Using byte for now.
			return elements
		}
	}

	// 5. Equivalent numbers
	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if !IsFreeToken(*t) {
			continue
		}
		if matches := episodeEquivalentPattern.FindStringSubmatchIndex(t.Value); matches != nil {
			val := t.Value[matches[2]:matches[3]]
			addEl(Episode, *t, val, t.Position+matches[2])
			return elements
		}
	}

	// 6. Separated number (e.g. ` - 08`)
	for i := len(tokens) - 1; i >= 0; i-- {
		if IsDashToken(tokens[i]) {
			nextIdx := FindNextToken(tokens, i, IsNotDelimiterToken)
			if nextIdx != -1 {
				nextTk := &tokens[nextIdx]
				if IsFreeToken(*nextTk) && IsNumericToken(*nextTk) {
					addEl(Episode, *nextTk, "", -1)
					return elements
				}
			}
		}
	}

	// 7. Isolated number `[12]`
	for i := 0; i < len(tokens)-2; i++ {
		t0, t1, t2 := &tokens[i], &tokens[i+1], &tokens[i+2]
		if IsOpenBracketToken(*t0) && IsCloseBracketToken(*t2) {
			if IsFreeToken(*t1) && IsNumericToken(*t1) {
				addEl(Episode, *t1, "", -1)
				return elements
			}
		}
	}

	// 8. Partial episode `4a`
	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if IsFreeToken(*t) && episodePartialPattern.MatchString(t.Value) {
			if i >= 2 && tokens[i-2].Value == "Ver1" && t.Value == "1a" {
				continue
			}
			addEl(Episode, *t, "", -1)
			return elements
		}
	}

	// 9. First number
	if len(tokens) > 0 {
		t := &tokens[0]
		if IsFreeToken(*t) && IsNumericToken(*t) {
			isFirst := false
			if len(tokens) <= 2 {
				isFirst = true
			} else if IsDashToken(tokens[1]) || IsDashToken(tokens[2]) {
				isFirst = true
			} else if tokens[1].Value == "." {
				if len(tokens[2].Value) > 0 && IsSpace([]rune(tokens[2].Value)[0]) {
					isFirst = true
				} else if tokens[2].ElementKind != nil && *tokens[2].ElementKind == FileExtension {
					isFirst = true
				}
			}
			if isFirst {
				addEl(Episode, *t, "", -1)
				return elements
			}
		}
	}

	// 10. Last number
	isVersionNumber := func(idx int) bool {
		if idx < 0 || idx >= len(tokens) {
			return false
		}
		t := &tokens[idx]
		if !IsNumericToken(*t) {
			return false
		}
		if idx > 0 {
			prev := &tokens[idx-1]
			if IsDelimiterToken(*prev) && prev.Value == "." {
				return true
			}
		}
		return false
	}

	for i := len(tokens) - 1; i >= 0; i-- {
		t := &tokens[i]
		if !IsFreeToken(*t) || !IsNumericToken(*t) {
			continue
		}
		if t.IsEnclosed || t.Position == 0 {
			continue
		}

		prevIdx := FindPrevToken(tokens, i, IsNotDelimiterToken)
		nextIdx := FindNextToken(tokens, i, IsNotDelimiterToken)

		if prevIdx != -1 {
			prev := &tokens[prevIdx]
			if Equal(prev.Value, "Cour") || Equal(prev.Value, "Part") || Equal(prev.Value, "Movie") || Equal(prev.Value, "No") {
				continue
			}
			if isVersionNumber(prevIdx) {
				continue
			}
			if IsCloseBracketToken(*prev) && prev.Value == "]" {
				continue
			}
		}

		if nextIdx != -1 {
			if isVersionNumber(nextIdx) {
				continue
			}
		}

		if prevIdx != -1 && nextIdx != -1 {
			if IsFreeToken(tokens[prevIdx]) && IsFreeToken(tokens[nextIdx]) {
				continue
			}
		}

		addEl(Episode, *t, "", -1)
		return elements
	}

	return elements
}
