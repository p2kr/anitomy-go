package anitomy

import (
	"strings"
)

type TokenKind int

const (
	TokenOpenBracket TokenKind = iota
	TokenCloseBracket
	TokenDelimiter
	TokenKeyword
	TokenText
)

type KeywordKind int

const (
	KeywordAudioChannels KeywordKind = iota
	KeywordAudioCodec
	KeywordAudioLanguage
	KeywordDevice
	KeywordEpisode
	KeywordEpisodeType
	KeywordLanguage
	KeywordOther
	KeywordPart
	KeywordReleaseGroup
	KeywordReleaseInformation
	KeywordReleaseVersion
	KeywordSeason
	KeywordSource
	KeywordSubtitles
	KeywordSubtitleLanguage
	KeywordType
	KeywordVideoCodec
	KeywordVideoColorDepth
	KeywordVideoDynamicRange
	KeywordVideoFormat
	KeywordVideoFrameRate
	KeywordVideoProfile
	KeywordVideoQuality
	KeywordVideoResolution
	KeywordVolume
)

type KeywordFlags uint8

const (
	FlagAmbiguous KeywordFlags = 1 << iota
	FlagSubword
	FlagPrefixForNumber
	FlagPrefixForOther
)

type Keyword struct {
	Kind  KeywordKind
	Flags KeywordFlags
}

func (k Keyword) IsAmbiguous() bool {
	return (k.Flags & FlagAmbiguous) != 0
}

func (k Keyword) IsSubword() bool {
	return (k.Flags & FlagSubword) != 0
}

func (k Keyword) IsPrefixForNumber() bool {
	return (k.Flags & FlagPrefixForNumber) != 0
}

func (k Keyword) IsPrefixForOther() bool {
	return (k.Flags & FlagPrefixForOther) != 0
}

type Token struct {
	Kind        TokenKind
	Value       string
	Keyword     *Keyword
	ElementKind *ElementKind
	Position    int  // index in input string (can be rune index or byte index, but typically string matches)
	IsEnclosed  bool // token is enclosed in brackets
	IsNumber    bool // all characters in `value` are digits
}

func IsIdentifiedToken(token Token) bool {
	return token.ElementKind != nil
}

func IsFreeToken(token Token) bool {
	return (token.Kind == TokenText || token.Kind == TokenKeyword) && token.ElementKind == nil
}

func IsOpenBracketToken(token Token) bool {
	return token.Kind == TokenOpenBracket
}

func IsCloseBracketToken(token Token) bool {
	return token.Kind == TokenCloseBracket
}

func IsBracketToken(token Token) bool {
	return IsOpenBracketToken(token) || IsCloseBracketToken(token)
}

func IsDashToken(token Token) bool {
	return token.Kind == TokenDelimiter && len(token.Value) > 0 && IsDash(rune(token.Value[0]))
}

func IsDelimiterToken(token Token) bool {
	return token.Kind == TokenDelimiter
}

func IsNotDelimiterToken(token Token) bool {
	return token.Kind != TokenDelimiter
}

func IsKeywordToken(token Token) bool {
	return token.Kind == TokenKeyword
}

func IsTextToken(token Token) bool {
	return token.Kind == TokenText
}

func IsEnclosedToken(token Token) bool {
	return token.IsEnclosed
}

func IsNumericToken(token Token) bool {
	return token.IsNumber
}

type entry struct {
	value string
	flags KeywordFlags
}

func makeBaseKeywords() map[KeywordKind][]entry {
	return map[KeywordKind][]entry{
		KeywordAudioChannels: {
			{"2.0", FlagAmbiguous},
			{"2.0ch", 0},
			{"2ch", 0},
			{"5.1", 0},
			{"5.1ch", 0},
			{"7.1", 0},
			{"7.1ch", 0},
		},
		KeywordAudioCodec: {
			{"AAC", FlagPrefixForOther},
			{"AACx2", 0},
			{"AACx3", 0},
			{"AACx4", 0},
			{"AC3", 0},
			{"EAC3", 0},
			{"E-AC-3", 0},
			{"E-AC3", 0},
			{"Atmos", 0},
			{"Dolby Atmos", 0},
			{"DD", FlagPrefixForOther},
			{"DDP", FlagPrefixForNumber},
			{"Dolby TrueHD", 0},
			{"TrueHD", FlagPrefixForNumber},
			{"DTS", FlagPrefixForNumber},
			{"DTS-ES", 0},
			{"FLAC", FlagPrefixForNumber},
			{"FLACx2", 0},
			{"FLACx3", 0},
			{"FLACx4", 0},
			{"Lossless", 0},
			{"MP3", 0},
			{"Opus", FlagAmbiguous},
			{"OGG", 0},
			{"Vorbis", 0},
		},
		KeywordAudioLanguage: {
			{"DualAudio", 0},
			{"Dual Audio", 0},
			{"MultiAudio", 0},
			{"Multi Audio", 0},
			{"Dub", 0},
			{"Dubbed", 0},
			{"Dubs", 0},
			{"ChiDub", 0},
			{"Chinese Dub", 0},
			{"EngDub", 0},
			{"English Dub", 0},
			{"GerDub", 0},
			{"German Dub", 0},
			{"JapDub", 0},
			{"Japanese Dub", 0},
			{"Korean Dub", 0},
		},
		KeywordDevice: {
			{"Android", FlagAmbiguous},
			{"iPad", FlagPrefixForNumber},
			{"iPhone", FlagPrefixForNumber},
			{"iPod", 0},
			{"PS", FlagPrefixForNumber},
			{"Xbox", FlagPrefixForNumber},
		},
		KeywordEpisode: {
			{"Ep", FlagPrefixForNumber},
			{"Eps", FlagPrefixForNumber},
			{"Episode", FlagPrefixForNumber},
			{"Episodes", FlagPrefixForNumber},
			{"Episodio", FlagPrefixForNumber},
			{"Episódio", FlagPrefixForNumber},
			{"Capitulo", FlagPrefixForNumber},
			{"Folge", FlagPrefixForNumber},
		},
		KeywordEpisodeType: {
			{"OP", FlagAmbiguous | FlagPrefixForNumber},
			{"Opening", FlagAmbiguous},
			{"NCOP", FlagPrefixForNumber},
			{"ED", FlagAmbiguous | FlagPrefixForNumber},
			{"Ending", FlagAmbiguous},
			{"NCED", FlagPrefixForNumber},
			{"Preview", FlagAmbiguous},
			{"PV", FlagAmbiguous | FlagPrefixForNumber},
		},
		KeywordLanguage: {
			{"CHS", 0},
			{"CHT", 0},
			{"ENG", 0},
			{"English", 0},
			{"ESP", FlagAmbiguous},
			{"Espanol", 0},
			{"Spanish", 0},
			{"ITA", FlagAmbiguous},
			{"JAP", 0},
			{"JPN", 0},
			{"PT-BR", 0},
			{"VOSTFR", 0},
		},
		KeywordOther: {
			{"Remaster", 0},
			{"Remastered", 0},
			{"Uncensored", 0},
			{"Uncut", 0},
			{"TS", 0},
			{"VFR", 0},
			{"Widescreen", 0},
			{"WS", 0},
		},
		KeywordPart: {
			{"Cour", FlagPrefixForNumber},
			{"Part", FlagAmbiguous | FlagPrefixForNumber},
			{"Parte", FlagPrefixForNumber},
		},
		KeywordReleaseGroup: {
			{"0x539", 0},
			{"THORA", 0},
			{"VARYG", 0},
		},
		KeywordReleaseInformation: {
			{"Batch", 0},
			{"Complete", 0},
			{"End", FlagAmbiguous},
			{"Final", FlagAmbiguous},
			{"Patch", 0},
			{"Remux", 0},
			{"Repack", 0},
		},
		KeywordReleaseVersion: {
			{"v0", 0},
			{"v1", 0},
			{"v2", 0},
			{"v3", 0},
			{"v4", 0},
		},
		KeywordSeason: {
			{"Season", FlagAmbiguous},
			{"Saison", FlagAmbiguous},
		},
		KeywordSource: {
			{"BD", FlagPrefixForOther},
			{"BDRip", 0},
			{"BluRay", 0},
			{"Blu ray", 0},
			{"DVD", FlagPrefixForNumber},
			{"DVD5", 0},
			{"DVD9", 0},
			{"DVDISO", 0},
			{"DVDRip", 0},
			{"DVD Rip", 0},
			{"R2DVD", 0},
			{"R2J", 0},
			{"R2JDVD", 0},
			{"R2JDVDRip", 0},
			{"HDTV", 0},
			{"HDTVRip", 0},
			{"TVRip", 0},
			{"TV Rip", 0},
			{"Web", FlagAmbiguous},
			{"Webcast", 0},
			{"WebDL", 0},
			{"Web DL", 0},
			{"WebRip", 0},
			{"ADN", 0},
			{"AMZN", 0},
			{"BILI", 0},
			{"Bilibili", 0},
			{"CR", 0},
			{"Crunchyroll", 0},
			{"DSNP", 0},
			{"Funi", 0},
			{"Funimation", 0},
			{"HIDI", 0},
			{"Hidive", 0},
			{"Hulu", 0},
			{"Netflix", 0},
			{"NF", 0},
			{"VRV", 0},
			{"YouTube", 0},
		},
		KeywordSubtitles: {
			{"ASS", 0},
			{"BIG5", 0},
			{"Hardsub", 0},
			{"Hardsubs", 0},
			{"RAW", 0},
			{"Softsub", 0},
			{"Softsubs", 0},
			{"Sub", 0},
			{"Subbed", 0},
			{"Subtitled", 0},
			{"Multisub", 0},
			{"Multi Sub", 0},
			{"Multi Subs", 0},
			{"Multiple Subtitle", 0},
		},
		KeywordSubtitleLanguage: {
			{"EngSub", 0},
			{"EngSubs", 0},
			{"GerSub", 0},
		},
		KeywordType: {
			{"TV", FlagAmbiguous},
			{"Movie", FlagAmbiguous},
			{"Gekijouban", FlagAmbiguous},
			{"OAD", FlagAmbiguous | FlagPrefixForNumber},
			{"OAV", FlagAmbiguous | FlagPrefixForNumber},
			{"OVA", FlagAmbiguous | FlagPrefixForNumber},
			{"ONA", FlagAmbiguous | FlagPrefixForNumber},
			{"SP", FlagAmbiguous | FlagPrefixForNumber},
			{"Special", FlagAmbiguous},
			{"Specials", FlagAmbiguous},
		},
		KeywordVideoColorDepth: {
			{"8bit", 0},
			{"8bits", 0},
			{"8 bit", 0},
			{"8 bits", 0},
			{"10bit", 0},
			{"10bits", 0},
			{"10 bit", 0},
			{"10 bits", 0},
		},
		KeywordVideoCodec: {
			{"AV1", 0},
			{"DivX", FlagPrefixForNumber},
			{"AVC", 0},
			{"H 264", 0},
			{"H264", 0},
			{"X 264", 0},
			{"X264", 0},
			{"H 265", 0},
			{"H265", 0},
			{"HEVC", FlagPrefixForNumber},
			{"X 265", 0},
			{"X265", 0},
			{"Xvid", 0},
		},
		KeywordVideoDynamicRange: {
			{"HDR", 0},
			{"HDR10", 0},
			{"Dolby Vision", 0},
			{"DV", 0},
		},
		KeywordVideoFormat: {
			{"AVI", 0},
			{"MP4", 0},
			{"RMVB", 0},
			{"WMV", 0},
			{"WMV3", 0},
			{"WMV9", 0},
		},
		KeywordVideoFrameRate: {
			{"23.976FPS", 0},
			{"24FPS", 0},
			{"29.97FPS", 0},
			{"30FPS", 0},
			{"60FPS", 0},
			{"120FPS", 0},
		},
		KeywordVideoProfile: {
			{"Hi10", 0},
			{"Hi10p", 0},
			{"Hi444", 0},
			{"Hi444P", 0},
			{"Hi444PP", 0},
		},
		KeywordVideoQuality: {
			{"HD", 0},
			{"SD", 0},
			{"HQ", 0},
			{"LQ", 0},
		},
		KeywordVideoResolution: {
			{"1080p", FlagSubword},
			{"1440p", FlagSubword},
			{"2160p", FlagSubword},
			{"4K", 0},
		},
		KeywordVolume: {
			{"Vol", FlagPrefixForNumber},
			{"Volume", FlagPrefixForNumber},
		},
	}
}

var Keywords map[string]Keyword

func init() {
	Keywords = make(map[string]Keyword)
	base := makeBaseKeywords()

	for kind, entries := range base {
		for _, e := range entries {
			kw := Keyword{Kind: kind, Flags: e.flags}
			key := strings.ToLower(e.value)
			Keywords[key] = kw
			if strings.Contains(e.value, " ") {
				for _, delimiter := range []string{"_", ".", "-"} {
					variant := strings.ReplaceAll(key, " ", delimiter)
					Keywords[variant] = kw
				}
			}
		}
	}
}
