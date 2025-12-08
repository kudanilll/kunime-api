package anime

import "context"

type Scraper interface {
    ScrapeAnimeDetail(ctx context.Context, slug string) (*Anime, error)
    ScrapeAnimeEpisodes(ctx context.Context, slug string) ([]Episode, error)
    SearchAnime(ctx context.Context, query string) ([]Anime, error)
	ScrapeOngoingAnime(ctx context.Context, page int) ([]OngoingAnime, error)
}

type Service struct {
	scraper Scraper
}

func NewService(scraper Scraper) *Service {
    return &Service{scraper: scraper}
}

func (s *Service) GetAnimeDetail(ctx context.Context, slug string) (*Anime, error) {
    // nanti kalau kamu mau cache di sini bisa
    return s.scraper.ScrapeAnimeDetail(ctx, slug)
}

func (s *Service) GetAnimeEpisodes(ctx context.Context, slug string) ([]Episode, error) {
    return s.scraper.ScrapeAnimeEpisodes(ctx, slug)
}

func (s *Service) SearchAnime(ctx context.Context, query string) ([]Anime, error) {
    return s.scraper.SearchAnime(ctx, query)
}

func (s *Service) GetOngoingAnime(ctx context.Context, page int) ([]OngoingAnime, error) {
	if page < 1 {
		page = 1
	}
	return s.scraper.ScrapeOngoingAnime(ctx, page)
}
