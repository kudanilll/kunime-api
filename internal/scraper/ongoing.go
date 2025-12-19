package scraper

import (
	"context"
	"fmt"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)


func (s *AnimeScraper) ScrapeOngoingAnime(ctx context.Context, page int) ([]anime.OngoingAnime, error) {
	ongoings := make([]anime.OngoingAnime, 0)
	c := newCollector(ctx, s.userAgent)

	c.OnHTML("div.venser div.venz ul li", func(e *colly.HTMLElement) {
		if err := ctx.Err(); err != nil {
			return
		}

		epText := strings.TrimSpace(e.ChildText(".epz"))      // "Episode 10"
		dayText := strings.TrimSpace(e.ChildText(".epztipe")) // "Sabtu"
		dateText := strings.TrimSpace(e.ChildText(".newnime"))// "06 Des"

		dayParts := strings.Fields(dayText)
		day := ""
		if len(dayParts) > 0 {
			day = dayParts[len(dayParts)-1]
		}

		href := strings.TrimSpace(e.ChildAttr(".thumb a", "href"))
		img := strings.TrimSpace(e.ChildAttr(".thumbz img", "src"))
		title := strings.TrimSpace(e.ChildText(".thumbz h2.jdlflm"))

		item := anime.OngoingAnime{
			Title:       title,
			Episode:     extractEpisodeNumber(epText),
			Day:         day,
			Date:        dateText,
			Image:		 absoluteURL(s.baseURL, img),
			Endpoint:    href,
		}

		ongoings = append(ongoings, item)
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	// mapping page -> URL:
	// page 1  => /ongoing-anime/
	// page 2+ => /ongoing-anime/page/{page}/
	var ongoingURL string
	if page <= 1 {
		ongoingURL = fmt.Sprintf("%s/ongoing-anime/", s.baseURL)
	} else {
		ongoingURL = fmt.Sprintf("%s/ongoing-anime/page/%d/", s.baseURL, page)
	}

	if err := c.Visit(ongoingURL); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	return ongoings, nil
}
