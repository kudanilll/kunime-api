package scraper

import (
	"context"
	"fmt"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

func (s *AnimeScraper) ScrapeCompletedAnime(ctx context.Context, page int) ([]anime.CompletedAnime, error) {
	acquire()
	defer release()

	completed := make([]anime.CompletedAnime, 0)
	c := newCollector(ctx, s.userAgent)

	c.OnHTML("div.rapi div.venz ul li", func(e *colly.HTMLElement) {
		if err := ctx.Err(); err != nil {
			return
		}

		// exam: "12 Episode"
		epText := strings.TrimSpace(e.ChildText(".epz"))
		episodes := extractEpisodeNumber(epText)

		// exam: "<i ...></i> 7.07" -> "7.07"
		scoreText := strings.TrimSpace(e.ChildText(".epztipe"))
		score := extractScore(scoreText)

		dateText := strings.TrimSpace(e.ChildText(".newnime"))
		href := strings.TrimSpace(e.ChildAttr(".thumb a", "href"))
		img := strings.TrimSpace(e.ChildAttr(".thumbz img", "src"))
		title := strings.TrimSpace(e.ChildText(".thumbz h2.jdlflm"))

		item := anime.CompletedAnime{
			Title:    title,
			Episodes: episodes,
			Score:    score,
			Date:     dateText,
			Image:    absoluteURL(s.baseURL, img),
			Endpoint: href,
		}

		completed = append(completed, item)
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	// page mapping:
	// 1  -> /complete-anime/
	// 2+ -> /complete-anime/page/{page}/
	var url string
	if page <= 1 {
		url = fmt.Sprintf("%s/complete-anime/", s.baseURL)
	} else {
		url = fmt.Sprintf("%s/complete-anime/page/%d/", s.baseURL, page)
	}

	if err := visitWithRetry(c, url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if len(completed) == 0 {
		return nil, fmt.Errorf("no completed anime found")
	}
	
	return completed, nil
}