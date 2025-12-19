package anime

import "context"

type Scraper interface {
	ScrapeOngoingAnime(ctx context.Context, page int) ([]OngoingAnime, error)
	ScrapeCompletedAnime(ctx context.Context, page int) ([]CompletedAnime, error)
	ScrapeGenres(ctx context.Context) ([]Genre, error)
	ScrapeGenrePage(ctx context.Context, slug string, page int) ([]GenrePageAnime, error)
	ScrapeAnimeBatch(ctx context.Context, slug string) (*AnimeBatch, error)
}

type Service struct {
	scraper Scraper
}

func NewService(scraper Scraper) *Service {
    return &Service{scraper: scraper}
}

func (s *Service) GetOngoingAnime(ctx context.Context, page int) ([]OngoingAnime, error) {
	if page < 1 {
		page = 1
	}
	return s.scraper.ScrapeOngoingAnime(ctx, page)
}

func (s *Service) GetCompletedAnime(ctx context.Context, page int) ([]CompletedAnime, error) {
	if page < 1 {
		page = 1
	}
	return s.scraper.ScrapeCompletedAnime(ctx, page)
}

func (s *Service) GetGenres(ctx context.Context) ([]Genre, error) {
	return s.scraper.ScrapeGenres(ctx)
}

func (s *Service) GetGenrePage(ctx context.Context, slug string, page int) ([]GenrePageAnime, error) {
	if page < 1 {
		page = 1
	}
	return s.scraper.ScrapeGenrePage(ctx, slug, page)
}

func (s *Service) GetAnimeBatch(ctx context.Context, slug string) (*AnimeBatch, error) {
	return s.scraper.ScrapeAnimeBatch(ctx, slug)
}
