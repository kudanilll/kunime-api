package scraper

import "context"

var scrapeLimiter = make(chan struct{}, 5)

func acquire(ctx context.Context) error {
	select {
	case scrapeLimiter <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func release() {
	<-scrapeLimiter
}
