package scraper

type AnimeScraper struct {
	baseURL   string
	userAgent string
}

func NewCollyScraper(baseURL string, userAgent string) *AnimeScraper {
	return &AnimeScraper{
		baseURL:   baseURL,
		userAgent: userAgent,
	}
}
