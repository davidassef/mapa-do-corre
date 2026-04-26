package handlers

import (
	"github.com/david/mapa-do-corre/backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type impactoService interface {
	ObterResumo() (models.ResumoImpacto, error)
}

type ImpactoHandler struct {
	service impactoService
}

func NewImpactoHandler(service impactoService) ImpactoHandler {
	return ImpactoHandler{service: service}
}

func (handler ImpactoHandler) Resumo(c *fiber.Ctx) error {
	resumo, err := handler.service.ObterResumo()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "falha ao carregar resumo de impacto",
		})
	}

	return c.JSON(fiber.Map{
		"resumo": resumo,
	})
}