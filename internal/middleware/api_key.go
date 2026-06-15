package middleware

import (
	"crypto/sha256"
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
)

func APIKeyMiddleware(validKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/healthz" {
			return c.Next()
		}

		key := c.Get("X-API-Key")
		if key == "" || !matchAPIKey(key, validKey) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid api key",
			})
		}
		return c.Next()
	}
}

// matchAPIKey compares two API keys in constant time using SHA-256 hashing
// to prevent timing attacks.
func matchAPIKey(provided, expected string) bool {
	hashedProvided := sha256.Sum256([]byte(provided))
	hashedExpected := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(hashedProvided[:], hashedExpected[:]) == 1
}
