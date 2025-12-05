package middleware

import "github.com/gofiber/fiber/v2"

func APIKeyMiddleware(validKey string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        key := c.Get("X-API-Key")
        if key == "" || key != validKey {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "invalid api key",
            })
        }
        return c.Next()
    }
}
