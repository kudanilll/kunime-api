package scraper

import (
	"time"

	"github.com/gocolly/colly/v2"
)

func visitWithRetry(c *colly.Collector, url string) error {
	// first attempt
	if err := c.Visit(url); err != nil {
		time.Sleep(300 * time.Millisecond)
		return c.Visit(url)
	}
	return nil
}
