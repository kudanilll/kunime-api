package scraper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kunime-api/internal/anime"
	"kunime-api/internal/config"

	"github.com/gocolly/colly/v2"
)

type AnimeScraper struct {
	baseURL string
}

func NewCollyScraper(baseURL string) *AnimeScraper {
	return &AnimeScraper{baseURL: baseURL}
}

func (s *AnimeScraper) ScrapeOngoingAnime(ctx context.Context, page int) ([]anime.OngoingAnime, error) {
	ongoings := make([]anime.OngoingAnime, 0)

	c := colly.NewCollector(
		colly.Async(false),
		colly.UserAgent(config.Load().UserAgent),
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

func (s *AnimeScraper) ScrapeCompletedAnime(ctx context.Context, page int) ([]anime.CompletedAnime, error) {
	completed := make([]anime.CompletedAnime, 0)

	c := colly.NewCollector(
		colly.Async(false),
		colly.UserAgent(config.Load().UserAgent),
	)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "otakudesu.*",
		RandomDelay: 800 * time.Millisecond,
	})

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			return
		}
	})

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

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	return completed, nil
}

func (s *AnimeScraper) ScrapeGenres(ctx context.Context) ([]anime.Genre, error) {
	genres := make([]anime.Genre, 0)

	c := colly.NewCollector(
		colly.Async(false),
	)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "otakudesu.*",
		RandomDelay: 800 * time.Millisecond,
	})

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			return
		}
	})

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

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	return genres, nil
}

// TODO: handle unknown episode
func (s *AnimeScraper) ScrapeGenrePage(ctx context.Context, slug string, page int) ([]anime.GenrePageAnime, error) {
	items := make([]anime.GenrePageAnime, 0)

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

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	return items, nil
}

func (s *AnimeScraper) ScrapeAnimeBatch(
	ctx context.Context,
	slug string,
) (*anime.AnimeBatch, error) {

	c := colly.NewCollector(
		colly.Async(false),
		colly.UserAgent(config.Load().UserAgent),
	)

	result := &anime.AnimeBatch{
		Qualities: make([]anime.BatchQuality, 0),
	}

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			return
		}
	})

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
	url := fmt.Sprintf("%s/batch/%s/", s.baseURL, slug)

	if err := c.Visit(url); err != nil {
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
