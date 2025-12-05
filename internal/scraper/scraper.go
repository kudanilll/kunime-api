package scraper

import (
	"context"
	"kunime-api/internal/anime"
)

type AnimeScraper struct {
	baseURL string
}

// ScrapeAnimeEpisodes implements anime.Scraper.
func (s *AnimeScraper) ScrapeAnimeEpisodes(ctx context.Context, slug string) ([]anime.Episode, error) {
	panic("unimplemented")
}

// SearchAnime implements anime.Scraper.
func (s *AnimeScraper) SearchAnime(ctx context.Context, query string) ([]anime.Anime, error) {
	panic("unimplemented")
}

func NewCollyScraper(baseURL string) *AnimeScraper {
	return &AnimeScraper{baseURL: baseURL}
}

func (s *AnimeScraper) ScrapeAnimeDetail(ctx context.Context, slug string) (*anime.Anime, error) {
	// implement pake colly di file colly_scraper.go
	// visit s.baseURL + "/anime/" + slug
	// parse HTML → isi struct anime.Anime
	panic("not implemented")
}
