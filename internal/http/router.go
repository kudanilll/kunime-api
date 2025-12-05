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

	app.Get("/", func (c *fiber.Ctx) error  {
		return c.JSON(fiber.Map{"message": "Welcome to Kunime API"})
	})

    api := app.Group("/api/v1")
    api.Get("/anime/:slug", h.GetAnimeDetail)
    api.Get("/anime/:slug/episodes", h.GetAnimeEpisodes)
    // api.Get("/search", h.SearchAnime)

    app.Get("/healthz", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })

    return app
}
