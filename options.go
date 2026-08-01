package anitomy

// Options configures which elements the Anitomy parser will attempt to extract.
// Setting a field to false will skip the corresponding extraction step, improving performance if that data is not needed.
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

// DefaultOptions returns an Options struct with all parsing features enabled,
// mirroring the default behavior of the C++ Anitomy library.
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
