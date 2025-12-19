package app

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
	once    sync.Once
	handler http.Handler
)

// inisitialize Fiber app
func initApp() {
	cfg := config.Load()
	scr := scraper.NewCollyScraper(
	    cfg.ScrapeBaseURL,
	    cfg.UserAgent,
	)
	svc := anime.NewService(scr)

	app := httpserver.NewServer(cfg, svc)

	handler = adaptor.FiberApp(app)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)

	if r.RequestURI == "" && r.URL != nil {
		r.RequestURI = r.URL.RequestURI()
	}

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
