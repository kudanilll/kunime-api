package scraper

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"kunime-api/internal/anime"

	"github.com/gocolly/colly/v2"
)

func (s *AnimeScraper) ScrapeAnimeEpisodes(
	ctx context.Context,
	animeSlug string,
) (*anime.AnimeEpisodeList, error) {

	c := newCollector(ctx, s.userAgent)

	episodes := make([]anime.AnimeEpisode, 0)

	c.OnHTML("div.episodelist ul li", func(e *colly.HTMLElement) {
		link := strings.TrimSpace(e.ChildAttr("a", "href"))
		title := strings.TrimSpace(e.ChildText("a"))

		if link == "" || title == "" {
			return
		}

		// extract episode number from title
		// example: "... Episode 12 (End) Subtitle Indonesia"
		epNum := extractEpisodeFromTitle(title)
		if epNum < 1 {
			return
		}

		episodes = append(episodes, anime.AnimeEpisode{
			Episode: epNum,
			URL:     link,
		})
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

	if len(episodes) == 0 {
		return nil, fmt.Errorf("episode list not found")
	}

	// IMPORTANT: sort from episode 1 → last
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].Episode < episodes[j].Episode
	})

	return &anime.AnimeEpisodeList{
		AnimeSlug: animeSlug,
		Episodes:  episodes,
	}, nil
}
