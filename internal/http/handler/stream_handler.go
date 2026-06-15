package handler

import (
	"context"
	"kunime-api/internal/anime"
	"time"

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

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetEpisodeStreams(ctx, slug)
	if err != nil {
		return respondError(c, err)
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

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	resolved, err := h.svc.ResolveStream(ctx, req.Token)
	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(fiber.Map{
    	"url": resolved.URL,
	})
}