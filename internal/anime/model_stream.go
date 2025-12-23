package anime

type StreamMirror struct {
	Quality string `json:"quality"` // 360p / 480p / 720p
	Server  string `json:"server"`  // otakuwatch5, mega, etc
	Token   string `json:"token"`   // data-content (base64)
}

type EpisodeStreams struct {
	EpisodeSlug string         `json:"episode_slug"`
	Streams     []StreamMirror `json:"streams"`
}

type ResolvedStream struct {
	Quality string `json:"quality"`
	Server  string `json:"server"`
	URL     string `json:"url"` // iframe src FINAL
}
