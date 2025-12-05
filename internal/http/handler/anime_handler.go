package handler

import (
	"kunime-api/internal/anime"

	"github.com/gofiber/fiber/v2"
)

type AnimeHandler struct {
    svc *anime.Service
}

func NewAnimeHandler(svc *anime.Service) *AnimeHandler {
    return &AnimeHandler{svc: svc}
}

func (h *AnimeHandler) GetAnimeDetail(c *fiber.Ctx) error {
    slug := c.Params("slug")

    result, err := h.svc.GetAnimeDetail(c.Context(), slug)
    if err != nil {
        return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(result)
}

func (h *AnimeHandler) GetAnimeEpisodes(c *fiber.Ctx) error {
    slug := c.Params("slug")

    episodes, err := h.svc.GetAnimeEpisodes(c.Context(), slug)
    if err != nil {
        return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{"items": episodes})
}
