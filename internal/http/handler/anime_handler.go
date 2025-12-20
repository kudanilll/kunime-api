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

	data, err := h.svc.GetGenrePage(c.UserContext(), slug, page)
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

	data, err := h.svc.GetAnimeBatch(c.Context(), animeSlug)
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

	data, err := h.svc.GetAnimeDetail(c.UserContext(), animeSlug)
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

	data, err := h.svc.GetAnimeEpisodes(c.UserContext(), animeSlug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}
