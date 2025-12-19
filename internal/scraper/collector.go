package scraper

import (
	"context"
	"time"

	"github.com/gocolly/colly/v2"
)

func newCollector(ctx context.Context, userAgent string) *colly.Collector {
	c := colly.NewCollector(
		colly.Async(false),
		colly.UserAgent(userAgent),
	)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 800 * time.Millisecond,
	})

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
		}
	})

	return c
}
