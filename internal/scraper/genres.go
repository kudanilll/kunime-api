package scraper

import (
	"context"
	"fmt"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

func (s *AnimeScraper) ScrapeGenres(ctx context.Context) ([]anime.Genre, error) {
	acquire()
	defer release()

	genres := make([]anime.Genre, 0)
	c := newCollector(ctx, s.userAgent)

	c.OnHTML("ul.genres li a", func(e *colly.HTMLElement) {
		if err := ctx.Err(); err != nil {
			return
		}

		name := strings.TrimSpace(e.Text)
		href := strings.TrimSpace(e.Attr("href"))

		if name == "" || href == "" {
			return
		}

		slug := extractGenreSlug(href)

		g := anime.Genre{
			Name:     name,
			Slug:     slug,
			Endpoint: absoluteURL(s.baseURL, href),
		}

		genres = append(genres, g)
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	url := fmt.Sprintf("%s/genre-list/", s.baseURL)
	if err := visitWithRetry(c, url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if len(genres) == 0 {
		return nil, fmt.Errorf("no genres found")
	}

	return genres, nil
}
