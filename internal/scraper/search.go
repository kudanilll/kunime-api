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

	c := newCollector(ctx, s.userAgent)

	results := make([]anime.AnimeSearchResult, 0)

	c.OnHTML(".page li", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("h2 > a"))
		href := strings.TrimSpace(e.ChildAttr("h2 > a", "href"))
		image := strings.TrimSpace(e.ChildAttr("img", "src"))
		meta := strings.TrimSpace(e.ChildText(".set"))

		if title == "" || href == "" {
			return
		}

		results = append(results, anime.AnimeSearchResult{
			Title:    title,
			Image:    absoluteURL(s.baseURL, image),
			Genres:   extractSearchGenres(meta),
			Status:   extractSearchStatus(meta),
			Rating:   extractSearchRating(meta),
			Endpoint: href,
		})
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	q := url.QueryEscape(query)
	searchURL := fmt.Sprintf(
		"%s/?s=%s&post_type=anime",
		s.baseURL,
		q,
	)

	if err := c.Visit(searchURL); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	return &anime.AnimeSearchResponse{
		Query: query,
		Data:  results,
	}, nil
}
