package anime

type OngoingAnime struct {
	Title    string `json:"title"`
	Episode  int    `json:"episode"`
	Day      string `json:"day"`  // Sabtu, Jumat, dll
	Date     string `json:"date"` // "06 Des"
	Image    string `json:"image"`
	Endpoint string `json:"endpoint"`
}

type CompletedAnime struct {
	Title    string  `json:"title"`
	Episodes int     `json:"episodes"`
	Score    float64 `json:"score"`
	Date     string  `json:"date"`
	Image    string  `json:"image"`
	Endpoint string  `json:"endpoint"`
}

type Genre struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Endpoint string `json:"endpoint"`
}

type GenrePageAnime struct {
	Title    string   `json:"title"`
	Endpoint string   `json:"endpoint"`
	Studio   string   `json:"studio"`
	Episodes string   `json:"episodes"` // "Unknown Eps", "12 Eps", "? Eps"
	Rating   string   `json:"rating"`   // must be string to accommodate "N/A"
	Genres   []string `json:"genres"`
	Image    string   `json:"image"`
	Synopsis string   `json:"synopsis"`
	Season   string   `json:"season"` // "Fall 2025"
}

type BatchLink struct {
	Server string `json:"server"` // OtakuDrive, Mega, etc
	URL    string `json:"url"`
}

type BatchQuality struct {
	Quality string      `json:"quality"` // MP4 360p, 480p, 720p
	Size    string      `json:"size"`    // 0.46 GB, 1.52 GB
	Links   []BatchLink `json:"links"`
}

type AnimeBatch struct {
	Title     string         `json:"title"`
	Qualities []BatchQuality `json:"qualities"`
}
