package anime

type OngoingAnime struct {
	Title       string `json:"title"`
	Episode     int    `json:"episode"`
	Day         string `json:"day"`   // Sabtu, Jumat, dll
	Date        string `json:"date"`  // "06 Des"
	Image       string `json:"image"`
    Endpoint    string `json:"endpoint"`
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
