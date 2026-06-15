package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getPageParam(c *fiber.Ctx) int {
	pageStr := c.Params("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// respondError logs the real error server-side and returns a generic
// error message to the client to prevent leaking internal details.
func respondError(c *fiber.Ctx, err error) error {
	log.Printf("[%s %s] error: %v", c.Method(), c.Path(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "internal server error",
	})
}
