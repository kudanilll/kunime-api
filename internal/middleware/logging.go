package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Format: [waktu] METHOD PATH STATUS LATENCY IP UA="user-agent" ERR=error
func Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// jalankan handler berikutnya
		err := c.Next()

		stop := time.Now()
		latency := stop.Sub(start)

		method := c.Method()
		path := c.OriginalURL()
		status := c.Response().StatusCode()
		ip := c.IP()
		ua := c.Get("User-Agent")

		log.Printf(
			"[%s] %s %s %d %s %s UA=%q ERR=%v",
			start.Format(time.RFC3339),
			method,
			path,
			status,
			latency,
			ip,
			ua,
			err,
		)

		return err
	}
}
