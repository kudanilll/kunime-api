package scraper

import (
	"context"
	"fmt"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

// TODO: handle unknown episode
func (s *AnimeScraper) ScrapeGenrePage(ctx context.Context, slug string, page int) ([]anime.GenrePageAnime, error) {
	acquire()
	defer release()

	items := make([]anime.GenrePageAnime, 0)
	c := newCollector(ctx, s.userAgent)

	c.OnHTML("div.venser div.col-anime-con", func(e *colly.HTMLElement) {
		if err := ctx.Err(); err != nil {
			return
		}

		title := strings.TrimSpace(e.ChildText(".col-anime-title a"))
		endpoint := strings.TrimSpace(e.ChildAttr(".col-anime-title a", "href"))
		studio := strings.TrimSpace(e.ChildText(".col-anime-studio"))
		eps := strings.TrimSpace(e.ChildText(".col-anime-eps"))
		rating := strings.TrimSpace(e.ChildText(".col-anime-rating"))
		image := strings.TrimSpace(e.ChildAttr(".col-anime-cover img", "src"))
		season := strings.TrimSpace(e.ChildText(".col-anime-date"))
		synopsis := strings.TrimSpace(e.ChildText(".col-synopsis"))

		// collect a list of genres
		var genres []string
		e.ForEach(".col-anime-genre a", func(_ int, g *colly.HTMLElement) {
			name := strings.TrimSpace(g.Text)
			if name != "" {
				genres = append(genres, name)
			}
		})

		if rating == "" {
			rating = "N/A"
		}

		item := anime.GenrePageAnime{
			Title:    title,
			Endpoint: endpoint,
			Studio:   studio,
			Episodes: eps,
			Rating:   rating,
			Genres:   genres,
			Image:    absoluteURL(s.baseURL, image),
			Synopsis: synopsis,
			Season:   season,
		}

		// at least have a title and endpoint; if empty, skip.
		if item.Title != "" && item.Endpoint != "" {
			items = append(items, item)
		}
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	url := fmt.Sprintf("%s/genres/%s/page/%d/", s.baseURL, slug, page)
	if err := visitWithRetry(c, url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no genre page anime found")
	}

	return items, nil
}
