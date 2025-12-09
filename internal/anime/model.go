package anime

type OngoingAnime struct {
	Title       string `json:"title"`
	Episode     int    `json:"episode"`
	Day         string `json:"day"`   // Sabtu, Jumat, dll
	Date        string `json:"date"`  // "06 Des"
	Image       string `json:"image"`
    Endpoint    string `json:"endpoint"`
}