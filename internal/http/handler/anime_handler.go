package handler

import (
	"context"
	"kunime-api/internal/anime"
	"strconv"
	"time"

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

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetOngoingAnime(ctx, page)

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

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetCompletedAnime(ctx, page)
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
	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetGenres(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

func (h *AnimeHandler) GetGenrePage(c *fiber.Ctx) error {
	slug := c.Params("genreSlug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "genre slug is required",
		})
	}

	pageStr := c.Params("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetGenrePage(ctx, slug, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"genre": slug,
		"page":  page,
		"data":  data,
	})
}

func (h *AnimeHandler) GetAnimeBatch(c *fiber.Ctx) error {
	animeSlug := c.Params("animeSlug")
	if animeSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "anime slug is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetAnimeBatch(ctx, animeSlug)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *AnimeHandler) GetAnimeDetail(c *fiber.Ctx) error {
	animeSlug := c.Params("animeSlug")
	if animeSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "anime slug is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetAnimeDetail(ctx, animeSlug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *AnimeHandler) GetAnimeEpisodes(c *fiber.Ctx) error {
	animeSlug := c.Params("animeSlug")
	if animeSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "anime slug is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.GetAnimeEpisodes(ctx, animeSlug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func (h *AnimeHandler) SearchAnime(c *fiber.Ctx) error {
	query := c.Params("query")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search query is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 15*time.Second)
	defer cancel()

	data, err := h.svc.SearchAnime(ctx, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}
