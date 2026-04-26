package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func NewPublicWriteLimiter(maxRequests int, windowSeconds int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: time.Duration(windowSeconds) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": fmt.Sprintf("limite de %d requisicoes em %d segundos excedido para este IP", maxRequests, windowSeconds),
			})
		},
	})
}
