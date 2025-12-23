package scraper

import (
	"context"
	"fmt"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

func (s *AnimeScraper) ScrapeEpisodeStreams(
	ctx context.Context,
	episodeSlug string,
) (*anime.EpisodeStreams, error) {
	acquire()
	defer release()

	c := newCollector(ctx, s.userAgent)
	results := make([]anime.StreamMirror, 0)

	c.OnHTML(".mirrorstream ul", func(e *colly.HTMLElement) {
		class := e.Attr("class") // m360p / m480p / m720p
		quality := strings.TrimPrefix(class, "m")
		
		if !strings.HasPrefix(class, "m") {
			return
		}

		e.ForEach("li a[data-content]", func(_ int, a *colly.HTMLElement) {
			server := strings.TrimSpace(a.Text)
			token := strings.TrimSpace(a.Attr("data-content"))

			if token == "" || server == "" {
				return
			}

			results = append(results, anime.StreamMirror{
				Quality: quality,
				Server:  server,
				Token:   token,
			})
		})
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	url := fmt.Sprintf("%s/episode/%s/", s.baseURL, episodeSlug)
	if err := visitWithRetry(c, url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no stream mirrors found")
	}

	return &anime.EpisodeStreams{
		EpisodeSlug: episodeSlug,
		Streams:     results,
	}, nil
}
