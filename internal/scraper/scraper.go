package scraper

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

type AnimeScraper struct {
	baseURL string
}

func NewCollyScraper(baseURL string) *AnimeScraper {
	return &AnimeScraper{baseURL: baseURL}
}

// ScrapeAnimeEpisodes implements anime.Scraper.
func (s *AnimeScraper) ScrapeAnimeEpisodes(ctx context.Context, slug string) ([]anime.Episode, error) {
	panic("unimplemented")
}

// SearchAnime implements anime.Scraper.
func (s *AnimeScraper) SearchAnime(ctx context.Context, query string) ([]anime.Anime, error) {
	panic("unimplemented")
}

func (s *AnimeScraper) ScrapeAnimeDetail(ctx context.Context, slug string) (*anime.Anime, error) {
	panic("not implemented")
}

func (s *AnimeScraper) ScrapeOngoingAnime(ctx context.Context, page int) ([]anime.OngoingAnime, error) {
	ongoings := make([]anime.OngoingAnime, 0)

	c := colly.NewCollector(
		colly.Async(false),
	)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 800 * time.Millisecond,
	})

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			return
		}
	})

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
			Immage:		 absoluteURL(s.baseURL, img),
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

// "Episode 10" -> 10
func extractEpisodeNumber(epText string) int {
	epText = strings.ToLower(epText)
	epText = strings.ReplaceAll(epText, "episode", "")
	epText = strings.TrimSpace(epText)

	if epText == "" {
		return 0
	}

	n, err := strconv.Atoi(epText)
	if err != nil {
		return 0
	}
	return n
}

func absoluteURL(base, p string) string {
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}

	u, err := url.Parse(base)
	if err != nil {
		return p
	}

	u.Path = path.Join(u.Path, p)
	return u.String()
}
