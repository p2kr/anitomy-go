package anitomy

import (
	"regexp"
)

var (
	seasonPattern   = regexp.MustCompile(`^S(\d{1,2})$`)
	seasonJpPattern = regexp.MustCompile(`^(?:第)?(\d{1,2})期$`)
)

func ParseSeason(tokens []Token) []Element {
	var elements []Element

	// Check ordinal/roman numbers around season keywords
	for i := 0; i < len(tokens)-2; i++ {
		t0, t1, t2 := &tokens[i], &tokens[i+1], &tokens[i+2]

		isSeasonKw := func(t Token) bool {
			return t.Keyword != nil && t.Keyword.Kind == KeywordSeason
		}

		// ends_with_season_keyword: free, delimiter, season_keyword
		if IsFreeToken(*t0) && IsDelimiterToken(*t1) && isSeasonKw(*t2) {
			if number := FromOrdinalNumber(t0.Value); number != "" {
				k := Season
				t0.ElementKind = &k
				t2.ElementKind = &k
				elements = append(elements, Element{Kind: Season, Value: number, Position: t0.Position})
				break
			}
		}

		// starts_with_season_keyword: season_keyword, delimiter, free
		if isSeasonKw(*t0) && IsDelimiterToken(*t1) && IsFreeToken(*t2) {
			var value string
			if IsNumericToken(*t2) {
				value = t2.Value
			} else if number := FromRomanNumber(t2.Value); number != "" {
				value = number
			}
			if value != "" {
				k := Season
				t0.ElementKind = &k
				t2.ElementKind = &k
				elements = append(elements, Element{Kind: Season, Value: value, Position: t2.Position})
				break
			}
		}
	}

	// Season pattern S2, S01-02
	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if !IsFreeToken(*t) {
			continue
		}

		if matches := seasonPattern.FindStringSubmatchIndex(t.Value); matches != nil {
			k := Season
			t.ElementKind = &k
			val := t.Value[matches[2]:matches[3]]
			elements = append(elements, Element{Kind: Season, Value: val, Position: t.Position + matches[2]})

			if i+2 < len(tokens) {
				next1, next2 := &tokens[i+1], &tokens[i+2]
				if IsDashToken(*next1) && IsFreeToken(*next2) && IsNumericToken(*next2) {
					next2.ElementKind = &k
					elements = append(elements, Element{Kind: Season, Value: next2.Value, Position: next2.Position})
					i += 2
				}
			}
		}
	}

	// Japanese counter
	if len(elements) == 0 {
		for i := 0; i < len(tokens); i++ {
			t := &tokens[i]
			if !IsFreeToken(*t) {
				continue
			}
			if matches := seasonJpPattern.FindStringSubmatchIndex(t.Value); matches != nil {
				k := Season
				t.ElementKind = &k
				val := t.Value[matches[2]:matches[3]]
				// Position requires rune offset matching. In go regex, position is byte offset.
				// For `第2期`, `第` is 3 bytes. `2` is at byte offset 3.
				// We need to convert byte offset to rune offset if Position is in runes,
				// but let's assume byte offsets for now or just adjust if needed.
				runeOffset := len([]rune(t.Value[:matches[2]]))
				elements = append(elements, Element{Kind: Season, Value: val, Position: t.Position + runeOffset})
				break
			}
		}
	}

	return elements
}
