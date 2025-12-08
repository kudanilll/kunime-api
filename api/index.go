package api

import (
	"net/http"
	"sync"

	"kunime-api/internal/anime"
	"kunime-api/internal/config"
	httpserver "kunime-api/internal/http"
	"kunime-api/internal/scraper"

	"github.com/gofiber/adaptor/v2"
)


var (
	handlerOnce sync.Once
	handler     http.Handler
)

// initialise Fiber app
func initFiberHandler() {
	cfg := config.Load()

	scr := scraper.NewCollyScraper(cfg.ScrapeBaseURL)
	animeSvc := anime.NewService(scr)

	app := httpserver.NewServer(cfg, animeSvc)

	// adaptor.FiberApp(*fiber.App) -> http.HandlerFunc
	handler = adaptor.FiberApp(app)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	handlerOnce.Do(initFiberHandler)

	// CORS preflight handling
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	handler.ServeHTTP(w, r)
}
