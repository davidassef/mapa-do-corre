package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCors(appEnv string, frontendOrigin string) fiber.Handler {
	allowedOrigins := buildAllowedOrigins(appEnv, frontendOrigin)

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	})
}

func buildAllowedOrigins(appEnv string, frontendOrigin string) []string {
	origins := []string{frontendOrigin}
	if strings.EqualFold(strings.TrimSpace(appEnv), "production") {
		return origins
	}

	return append(origins,
		"http://127.0.0.1:4173",
		"http://127.0.0.1:5173",
		"http://localhost:4173",
		"http://localhost:5173",
	)
}
