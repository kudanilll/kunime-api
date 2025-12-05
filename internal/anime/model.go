package anime

type Anime struct {
    Title       string   `json:"title"`
    Slug        string   `json:"slug"`
    Synopsis    string   `json:"synopsis"`
    PosterImage string   `json:"posterImage"`
    Genres      []string `json:"genres"`
    Status      string   `json:"status"`
    Year        int      `json:"year"`
    Season      string   `json:"season"`
}

type Episode struct {
    Number    int    `json:"number"`
    Title     string `json:"title"`
    Synopsis  string `json:"synopsis"`
    VideoURL  string `json:"videoUrl"` // link video yg kamu kirim ke client
}
