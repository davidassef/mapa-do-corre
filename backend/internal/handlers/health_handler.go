package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type healthChecker interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	database healthChecker
}

func NewHealthHandler(database healthChecker) HealthHandler {
	return HealthHandler{database: database}
}

func (handler HealthHandler) Handle(c *fiber.Ctx) error {
	response := fiber.Map{
		"status": "up",
		"checks": fiber.Map{
			"api": fiber.Map{
				"status": "up",
			},
			"database": fiber.Map{
				"status": "not_configured",
			},
		},
	}

	if handler.database == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(response)
	}

	healthCtx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	err := handler.database.Ping(healthCtx)
	if err != nil {
		response["status"] = "degraded"
		response["checks"].(fiber.Map)["database"] = fiber.Map{
			"status":  "down",
			"message": err.Error(),
		}

		return c.Status(fiber.StatusServiceUnavailable).JSON(response)
	}

	response["checks"].(fiber.Map)["database"] = fiber.Map{
		"status": "up",
	}

	return c.JSON(response)
}