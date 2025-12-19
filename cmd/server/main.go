package main

import (
	"log"

	"github.com/joho/godotenv"

	"kunime-api/internal/anime"
	"kunime-api/internal/config"
	"kunime-api/internal/http"
	"kunime-api/internal/scraper"
)

func main() {

    // load .env
    if err := godotenv.Load(); err != nil {
        log.Println("no .env file found (or failed to load), using system env")
    }

    cfg := config.Load()
    scr := scraper.NewCollyScraper(
	    cfg.ScrapeBaseURL,
	    cfg.UserAgent,
    )
    animeService := anime.NewService(scr)
    app := http.NewServer(cfg, animeService)

    if err := app.Listen(":" + cfg.Port); err != nil {
        log.Fatal(err)
    }
}
