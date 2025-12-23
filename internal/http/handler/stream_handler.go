package handler

import (
	"kunime-api/internal/anime"

	"github.com/gofiber/fiber/v2"
)

type StreamHandler struct {
	svc *anime.Service
}

func NewStreamHandler(svc *anime.Service) *StreamHandler {
	return &StreamHandler{svc: svc}
}

func (h *StreamHandler) GetEpisodeStreams(c *fiber.Ctx) error {
	slug := c.Params("episodeSlug")
	data, err := h.svc.GetEpisodeStreams(c.Context(), slug)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (h *StreamHandler) ResolveStream(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil || req.Token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "token required"})
	}

	resolved, err := h.svc.ResolveStream(c.Context(), req.Token)
	if err != nil {
    	return c.Status(500).JSON(fiber.Map{
        	"error": err.Error(),
    	})
	}

	return c.JSON(fiber.Map{
    	"url": resolved.URL,
	})
}