package scraper

import (
	"net/http"
	"time"
)

type AnimeScraper struct {
	baseURL   string
	userAgent string
	client    *http.Client
}

func NewCollyScraper(baseURL string, userAgent string) *AnimeScraper {
	return &AnimeScraper{
		baseURL:   baseURL,
		userAgent: userAgent,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
