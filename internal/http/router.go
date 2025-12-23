package http

import (
	"kunime-api/internal/anime"
	"kunime-api/internal/config"
	"kunime-api/internal/http/handler"
	"kunime-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewServer(cfg config.Config, animeSvc *anime.Service) *fiber.App {
    app := fiber.New()

    // middleware global
    app.Use(middleware.Logging())
    app.Use(middleware.APIKeyMiddleware(cfg.APIKey))

    h := handler.NewAnimeHandler(animeSvc)
    streamH := handler.NewStreamHandler(animeSvc)

	app.Get("/", func (c *fiber.Ctx) error  {
        return c.JSON(fiber.Map {
            "github": "https://github.com/kudanilll/kunime-api", 
            "support": "https://buymeacoffee.com/kudanil",
            "documentation": "https://github.com/kudanilll/kunime-api/blob/master/docs/API.md",
            "endpoint": fiber.Map {
                "Get Ongoing Anime":         "/api/v1/ongoing-anime/:page", 
                "Get Completed Anime":       "/api/v1/completed-anime/:page", 
                "Get Genres":                "/api/v1/genres",
				"Get Anime by Genre & Page": "/api/v1/genre/:genreSlug/:page",
                "Get Anime Batch":           "/api/v1/anime/:animeSlug/batch",
                "Get Anime Detail":          "/api/v1/anime/:animeSlug",
                "Get Anime Episodes":        "/api/v1/anime/:animeSlug/episodes",
                "Search Anime":              "/api/v1/search/:query",
                "Get Anime Streams":         "/api/v1/anime/:episodeSlug/streams",
                "Resolve Stream":            "/api/v1/streams/resolve - POST",
            },
        })
	})

    api := app.Group("/api/v1")

    // ongoing anime
    api.Get("/ongoing-anime", h.GetOngoingAnime)
    api.Get("/ongoing-anime/:page", h.GetOngoingAnime)

	// completed anime
	api.Get("/completed-anime", h.GetCompletedAnime)
	api.Get("/completed-anime/:page", h.GetCompletedAnime)

	// genres
	api.Get("/genres", h.GetGenres)
	
    // genre page
	api.Get("/genre/:genreSlug/:page", h.GetGenrePage)

    // anime batch, detail, episode list
    api.Get("/anime/:animeSlug/batch", h.GetAnimeBatch)
    api.Get("/anime/:animeSlug", h.GetAnimeDetail)
    api.Get("/anime/:animeSlug/episodes", h.GetAnimeEpisodes)

    // search
    api.Get("/search/:query", h.SearchAnime)

    // stream endpoints
    api.Get("/anime/:episodeSlug/streams", streamH.GetEpisodeStreams)
    api.Post("/streams/resolve", streamH.ResolveStream)
    
    // health check
    app.Get("/healthz", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })

    return app
}
