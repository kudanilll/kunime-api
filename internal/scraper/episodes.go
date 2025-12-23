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
	acquire()
	defer release()

	c := newCollector(ctx, s.userAgent)
	episodes := make([]anime.AnimeEpisode, 0)

	c.OnHTML("div.episodelist ul > li", func(e *colly.HTMLElement) {
		href := strings.TrimSpace(e.ChildAttr("a", "href"))
		title := strings.TrimSpace(e.ChildText("a"))

		if href == "" || title == "" {
			return
		}

		// skip batch
		if !strings.Contains(href, "/episode/") {
			return
		}

		// extract episode number
		epNum := extractEpisodeFromTitle(title)
		if epNum < 1 {
			return
		}

		// convert full URL → slug only
		// https://otakudesu.best/episode/kni-s2-episode-12-sub-indo/
		// → kni-s2-episode-12-sub-indo
		epSlug := extractEpisodeSlug(href)
		if epSlug == "" {
			return
		}

		episodes = append(episodes, anime.AnimeEpisode{
			Episode: epNum,
			Slug:    epSlug,
		})
	})

	var scrapeErr error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = err
	})

	url := fmt.Sprintf("%s/anime/%s/", s.baseURL, animeSlug)
	if err := visitWithRetry(c, url); err != nil {
		return nil, err
	}

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	if len(episodes) == 0 {
		return nil, fmt.Errorf("episode list not found")
	}

	// sort ascending
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].Episode < episodes[j].Episode
	})

	return &anime.AnimeEpisodeList{
		AnimeSlug: animeSlug,
		Episodes:  episodes,
	}, nil
}
