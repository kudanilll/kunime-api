package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

func (s *AnimeScraper) SearchAnime(
	ctx context.Context,
	query string,
) (*anime.AnimeSearchResponse, error) {

	if err := acquire(ctx); err != nil {
		return nil, err
	}
	defer release()

	c := newCollector(ctx, s.userAgent)
	results := make([]anime.AnimeSearchResult, 0)

	c.OnHTML(".page ul.chivsrc > li", func(e *colly.HTMLElement) {
		rawTitle := strings.TrimSpace(e.ChildText("h2 > a"))
		title := cleanAnimeTitle(rawTitle)

		href := strings.TrimSpace(e.ChildAttr("h2 > a", "href"))
		image := strings.TrimSpace(e.ChildAttr("img", "src"))

		if title == "" || href == "" {
			return
		}

		var (
			genres []string
			status string
			rating string
		)

		e.ForEach(".set", func(_ int, s *colly.HTMLElement) {
			label := strings.ToLower(strings.TrimSpace(s.ChildText("b")))
			value := strings.TrimSpace(strings.ReplaceAll(s.Text, s.ChildText("b"), ""))

			switch {
			case strings.Contains(label, "genres"):
				s.ForEach("a", func(_ int, a *colly.HTMLElement) {
					g := strings.TrimSpace(a.Text)
					if g != "" {
						genres = append(genres, g)
					}
				})

			case strings.Contains(label, "status"):
				status = strings.TrimPrefix(value, ":")

			case strings.Contains(label, "rating"):
				rating = strings.TrimPrefix(value, ":")
			}
		})

		results = append(results, anime.AnimeSearchResult{
			Title:    title,
			Image:    absoluteURL(s.baseURL, image),
			Genres:   genres,
			Status:   strings.TrimSpace(status),
			Rating:   strings.TrimSpace(rating),
			Endpoint: href,
		})
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	searchURL := fmt.Sprintf(
		"%s/?s=%s&post_type=anime",
		s.baseURL,
		url.QueryEscape(query),
	)

	if err := visitWithRetry(c, searchURL); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no search results found")
	}

	return &anime.AnimeSearchResponse{
		Query: query,
		Data:  results,
	}, nil
}
