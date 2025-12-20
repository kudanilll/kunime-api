package scraper

import (
	"context"
	"fmt"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

func (s *AnimeScraper) ScrapeAnimeDetail(
	ctx context.Context,
	animeSlug string,
) (*anime.AnimeDetail, error) {

	c := newCollector(ctx, s.userAgent)

	result := &anime.AnimeDetail{
		Genres:    make([]string, 0),
		Producers: make([]string, 0),
	}

	c.OnHTML("div.fotoanime", func(e *colly.HTMLElement) {
		result.Image = absoluteURL(
			s.baseURL,
			strings.TrimSpace(e.ChildAttr("img", "src")),
		)

		e.ForEach("div.infozingle p span", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)

			switch {
			case strings.HasPrefix(text, "Judul"):
				result.Title = extractValue(text)
			case strings.HasPrefix(text, "Japanese"):
				result.JapaneseTitle = extractValue(text)
			case strings.HasPrefix(text, "Skor"):
				result.Score = extractValue(text)
			case strings.HasPrefix(text, "Tipe"):
				result.Type = extractValue(text)
			case strings.HasPrefix(text, "Status"):
				result.Status = extractValue(text)
			case strings.HasPrefix(text, "Total Episode"):
				result.TotalEpisode = extractValue(text)
			case strings.HasPrefix(text, "Durasi"):
				result.Duration = extractValue(text)
			case strings.HasPrefix(text, "Tanggal Rilis"):
				result.ReleaseDate = extractValue(text)
			case strings.HasPrefix(text, "Studio"):
				result.Studio = extractValue(text)
			case strings.HasPrefix(text, "Produser"):
				raw := extractValue(text)
				for _, p := range strings.Split(raw, ",") {
					p = strings.TrimSpace(p)
					if p != "" {
						result.Producers = append(result.Producers, p)
					}
				}
			}
		})

		e.ForEach("div.infozingle a[rel='tag']", func(_ int, g *colly.HTMLElement) {
			genre := strings.TrimSpace(g.Text)
			if genre != "" {
				result.Genres = append(result.Genres, genre)
			}
		})
	})

	c.OnHTML("div.sinopc", func(e *colly.HTMLElement) {
		result.Synopsis = strings.TrimSpace(e.Text)
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	url := fmt.Sprintf("%s/anime/%s/", s.baseURL, animeSlug)
	if err := c.Visit(url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if result.Title == "" {
		return nil, fmt.Errorf("anime detail not found")
	}

	return result, nil
}
