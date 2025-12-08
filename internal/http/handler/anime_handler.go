package handler

import (
	"kunime-api/internal/anime"
	"strconv"

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

func (h *AnimeHandler) GetOngoingAnime(c *fiber.Ctx) error {
	pageStr := c.Params("page", "1")
	page, err := strconv.Atoi(pageStr)
    
	if err != nil || page < 1 {
		page = 1
	}

	items, err := h.svc.GetOngoingAnime(c.Context(), page)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"page":  page,
		"items": items,
	})
}
