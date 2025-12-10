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

func (h *AnimeHandler) GetOngoingAnime(c *fiber.Ctx) error {
	page := getPageParam(c)

	data, err := h.svc.GetOngoingAnime(c.UserContext(), page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"page": page,
		"data": data,
	})
}

func (h *AnimeHandler) GetCompletedAnime(c *fiber.Ctx) error {
	page := getPageParam(c)

	data, err := h.svc.GetCompletedAnime(c.UserContext(), page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"page": page,
		"data": data,
	})
}

func (h *AnimeHandler) GetGenres(c *fiber.Ctx) error {
	data, err := h.svc.GetGenres(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
