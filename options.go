package anitomy

type Options struct {
	ParseEpisode         bool
	ParseEpisodeTitle    bool
	ParseFileChecksum    bool
	ParseFileExtension   bool
	ParsePart            bool
	ParseReleaseGroup    bool
	ParseSeason          bool
	ParseTitle           bool
	ParseVideoResolution bool
	ParseYear            bool
}

func DefaultOptions() Options {
	return Options{
		ParseEpisode:         true,
		ParseEpisodeTitle:    true,
		ParseFileChecksum:    true,
		ParseFileExtension:   true,
		ParsePart:            true,
		ParseReleaseGroup:    true,
		ParseSeason:          true,
		ParseTitle:           true,
		ParseVideoResolution: true,
		ParseYear:            true,
	}
}
