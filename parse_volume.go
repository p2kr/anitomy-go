package anitomy

import (
	"regexp"
)

var (
	volumeSinglePattern   = regexp.MustCompile(`^(\d{1,4})(?:[vV](\d))?$`)
	volumeMultiplePattern = regexp.MustCompile(`^(\d{1,4})&(\d{1,4})$`)
)

func ParseVolume(tokens []Token) []Element {
	var elements []Element

	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if t.Keyword != nil && t.Keyword.Kind == KeywordVolume {
			nextIdx := FindNextToken(tokens, i, IsNotDelimiterToken)
			if nextIdx == -1 || !IsFreeToken(tokens[nextIdx]) {
				continue
			}
			nextToken := &tokens[nextIdx]

			if matches := volumeSinglePattern.FindStringSubmatchIndex(nextToken.Value); matches != nil {
				k := Volume
				t.ElementKind = &k
				nextToken.ElementKind = &k

				volVal := nextToken.Value[matches[2]:matches[3]]
				elements = append(elements, ElementFromToken(Volume, *nextToken, volVal, -1))

				if matches[4] != -1 && matches[5] != -1 {
					verVal := nextToken.Value[matches[4]:matches[5]]
					// position is nextToken.Position + char offset. Since we use bytes for index in Go strings, it's fine for ascii digits.
					pos := nextToken.Position + matches[4]
					elements = append(elements, Element{Kind: ReleaseVersion, Value: verVal, Position: pos})
				}
			} else if matches := volumeMultiplePattern.FindStringSubmatchIndex(nextToken.Value); matches != nil {
				k := Volume
				t.ElementKind = &k
				nextToken.ElementKind = &k

				vol1 := nextToken.Value[matches[2]:matches[3]]
				vol2 := nextToken.Value[matches[4]:matches[5]]
				elements = append(elements, ElementFromToken(Volume, *nextToken, vol1, -1))
				elements = append(elements, ElementFromToken(Volume, *nextToken, vol2, -1))
			}

			i = nextIdx
		}
	}
	return elements
}
