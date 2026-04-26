package handlers

import (
	"strconv"
	"strings"

	"github.com/david/mapa-do-corre/backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type correService interface {
	BuscarProximos(filtro models.FiltroBuscaCorre) ([]models.PrestadorResumo, error)
}

type CorreHandler struct {
	service correService
}

func NewCorreHandler(service correService) CorreHandler {
	return CorreHandler{service: service}
}

func (handler CorreHandler) Buscar(c *fiber.Ctx) error {
	filtro, err := montarFiltroBusca(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	prestadores, err := handler.service.BuscarProximos(filtro)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "falha ao buscar corres proximos",
		})
	}

	return c.JSON(fiber.Map{
		"prestadores": prestadores,
		"meta": fiber.Map{
			"quantidade":  len(prestadores),
			"raioMetros":  filtro.RaioMetros,
			"categoria":   filtro.Categoria,
			"latitude":    filtro.Latitude,
			"longitude":   filtro.Longitude,
			"implementado": true,
		},
	})
}

func montarFiltroBusca(c *fiber.Ctx) (models.FiltroBuscaCorre, error) {
	latitude, err := parseRequiredFloat(c.Query("lat"), "lat")
	if err != nil {
		return models.FiltroBuscaCorre{}, err
	}

	longitude, err := parseRequiredFloat(c.Query("lon"), "lon")
	if err != nil {
		return models.FiltroBuscaCorre{}, err
	}

	raioMetros := 1500
	raioBruto := strings.TrimSpace(c.Query("raioMetros"))
	if raioBruto != "" {
		parsedRaio, err := strconv.Atoi(raioBruto)
		if err != nil {
			return models.FiltroBuscaCorre{}, fiber.NewError(fiber.StatusBadRequest, "raioMetros precisa ser numerico")
		}

		raioMetros = parsedRaio
	}

	return models.FiltroBuscaCorre{
		Latitude:   latitude,
		Longitude:  longitude,
		RaioMetros: raioMetros,
		Categoria:  strings.TrimSpace(c.Query("categoria")),
	}, nil
}

func parseRequiredFloat(rawValue string, fieldName string) (float64, error) {
	value := strings.TrimSpace(rawValue)
	if value == "" {
		return 0, fiber.NewError(fiber.StatusBadRequest, fieldName+" e obrigatorio")
	}

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, fieldName+" precisa ser numerico")
	}

	return parsedValue, nil
}