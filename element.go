package anitomy

type ElementKind string

const (
	AudioTerm          ElementKind = "audio_term"
	Device             ElementKind = "device"
	Episode            ElementKind = "episode"
	EpisodeTitle       ElementKind = "episode_title"
	FileChecksum       ElementKind = "file_checksum"
	FileExtension      ElementKind = "file_extension"
	Language           ElementKind = "language"
	Other              ElementKind = "other"
	Part               ElementKind = "part"
	ReleaseGroup       ElementKind = "release_group"
	ReleaseInformation ElementKind = "release_information"
	ReleaseVersion     ElementKind = "release_version"
	Season             ElementKind = "season"
	Source             ElementKind = "source"
	Subtitles          ElementKind = "subtitles"
	Title              ElementKind = "title"
	Type               ElementKind = "type"
	VideoResolution    ElementKind = "video_resolution"
	VideoTerm          ElementKind = "video_term"
	Volume             ElementKind = "volume"
	Year               ElementKind = "year"
)

type Element struct {
	Kind     ElementKind
	Value    string
	Position int
}
