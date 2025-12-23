package anime

import "context"

func (s *Service) GetEpisodeStreams(
	ctx context.Context,
	episodeSlug string,
) (*EpisodeStreams, error) {
	return s.scraper.ScrapeEpisodeStreams(ctx, episodeSlug)
}

func (s *Service) ResolveStream(
	ctx context.Context,
	token string,
) (*ResolvedStream, error) {
	url, err := s.scraper.ResolveStreamURL(ctx, token)
	if err != nil {
		return nil, err
	}

	return &ResolvedStream{
		URL: url,
	}, nil
}
