package scraper

import (
	"context"
	"fmt"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

func (s *AnimeScraper) ScrapeAnimeBatch(
	ctx context.Context,
	animeSlug string,
) (*anime.AnimeBatch, error) {
	acquire()
	defer release()

	c := newCollector(ctx, s.userAgent)
	result := &anime.AnimeBatch{
		Qualities: make([]anime.BatchQuality, 0),
	}

	c.OnHTML("div.batchlink", func(e *colly.HTMLElement) {
		result.Title = strings.TrimSpace(e.ChildText("h4"))

		e.ForEach("ul > li", func(_ int, li *colly.HTMLElement) {
			quality := strings.TrimSpace(li.ChildText("strong"))
			size := strings.TrimSpace(li.ChildText("i"))

			if quality == "" {
				return
			}

			q := anime.BatchQuality{
				Quality: quality,
				Size:    size,
				Links:   make([]anime.BatchLink, 0),
			}

			li.ForEach("a", func(_ int, a *colly.HTMLElement) {
				server := strings.TrimSpace(a.Text)
				url := strings.TrimSpace(a.Attr("href"))

				if server == "" || url == "" {
					return
				}

				q.Links = append(q.Links, anime.BatchLink{
					Server: server,
					URL:    url,
				})
			})

			if len(q.Links) > 0 {
				result.Qualities = append(result.Qualities, q)
			}
		})
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	// IMPORTANT: batch endpoint, not anime endpoint
	url := fmt.Sprintf("%s/batch/%s/", s.baseURL, animeSlug)

	if err := visitWithRetry(c, url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if result.Title == "" || len(result.Qualities) == 0 {
		return nil, fmt.Errorf("batch not found or empty")
	}

	return result, nil
}
