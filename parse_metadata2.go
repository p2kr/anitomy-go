package anitomy

func ParseYear(tokens []Token) *Element {
	for i := 0; i < len(tokens)-2; i++ {
		t0, t1, t2 := tokens[i], &tokens[i+1], tokens[i+2]

		if IsOpenBracketToken(t0) && IsCloseBracketToken(t2) && IsFreeToken(*t1) && IsNumericToken(*t1) {
			num := ToInt(t1.Value)
			if 1950 < num && num < 2050 {
				k := Year
				t1.ElementKind = &k
				e := ElementFromToken(Year, *t1, "", -1)
				return &e
			}
		}
	}
	return nil
}

func ParsePart(tokens []Token) *Element {
	for i := 0; i < len(tokens); i++ {
		t := &tokens[i]
		if t.Keyword != nil && t.Keyword.Kind == KeywordPart {
			nextIdx := FindNextToken(tokens, i, IsNotDelimiterToken)
			if nextIdx != -1 && IsNumericToken(tokens[nextIdx]) {
				k := Part
				t.ElementKind = &k
				tokens[nextIdx].ElementKind = &k
				e := ElementFromToken(Part, tokens[nextIdx], "", -1)
				return &e
			}
		}
	}
	return nil
}

func ParseVideoResolution(tokens []Token) []Element {
	return nil // stub for now
}

func ParseKeywords(tokens []Token, options Options) []Element {
	var elements []Element
	for i := range tokens {
		t := &tokens[i]
		if t.Keyword != nil && !IsIdentifiedToken(*t) {
			if t.Keyword.IsAmbiguous() {
				continue
			}
			// Mapping from KeywordKind to ElementKind
			var ek ElementKind
			switch t.Keyword.Kind {
			case KeywordAudioChannels, KeywordAudioCodec, KeywordAudioLanguage:
				ek = AudioTerm
			case KeywordDevice:
				ek = Device
			case KeywordEpisode:
				continue // handled in parse_episode
			case KeywordEpisodeType:
				ek = Type
			case KeywordLanguage:
				ek = Language
			case KeywordOther:
				ek = Other
			case KeywordPart:
				continue
			case KeywordReleaseGroup:
				ek = ReleaseGroup
			case KeywordReleaseInformation:
				ek = ReleaseInformation
			case KeywordReleaseVersion:
				ek = ReleaseVersion
			case KeywordSeason:
				continue
			case KeywordSource:
				ek = Source
			case KeywordSubtitles, KeywordSubtitleLanguage:
				ek = Subtitles
			case KeywordType:
				ek = Type
			case KeywordVideoCodec, KeywordVideoColorDepth, KeywordVideoDynamicRange, KeywordVideoFormat, KeywordVideoFrameRate, KeywordVideoProfile, KeywordVideoQuality:
				ek = VideoTerm
			case KeywordVideoResolution:
				ek = VideoResolution
			case KeywordVolume:
				continue
			}

			if ek != "" {
				t.ElementKind = &ek
				elements = append(elements, ElementFromToken(ek, *t, "", -1))
			}
		}
	}
	return elements
}
