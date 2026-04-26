package service

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/david/mapa-do-corre/backend/internal/models"
)

type correFinder interface {
	BuscarProximos(ctx context.Context, filtro models.FiltroBuscaCorre) ([]models.PrestadorResumo, error)
}

type CorreService struct {
	store      *CorreStore
	repository correFinder
}

func NewCorreService(store *CorreStore, repository correFinder) *CorreService {
	return &CorreService{store: store, repository: repository}
}

func (service *CorreService) BuscarProximos(filtro models.FiltroBuscaCorre) ([]models.PrestadorResumo, error) {
	filtro = normalizarFiltroBusca(filtro)

	if service.repository != nil {
		queryCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return service.repository.BuscarProximos(queryCtx, filtro)
	}

	prestadores := service.store.ListarPrestadores()
	if len(prestadores) == 0 {
		return []models.PrestadorResumo{}, nil
	}

	categoriaFiltro := strings.ToLower(strings.TrimSpace(filtro.Categoria))
	prestadoresProximos := make([]models.PrestadorResumo, 0, len(prestadores))

	for _, prestador := range prestadores {
		if categoriaFiltro != "" && strings.ToLower(prestador.Categoria) != categoriaFiltro {
			continue
		}

		distanciaMetros := calcularDistanciaMetros(filtro.Latitude, filtro.Longitude, prestador.Latitude, prestador.Longitude)
		if int(distanciaMetros) > filtro.RaioMetros {
			continue
		}

		prestador.DistanciaMetros = int(math.Round(distanciaMetros))
		prestadoresProximos = append(prestadoresProximos, prestador)
	}

	sort.Slice(prestadoresProximos, func(leftIndex int, rightIndex int) bool {
		return prestadoresProximos[leftIndex].DistanciaMetros < prestadoresProximos[rightIndex].DistanciaMetros
	})

	return prestadoresProximos, nil
}

func normalizarFiltroBusca(filtro models.FiltroBuscaCorre) models.FiltroBuscaCorre {
	if filtro.RaioMetros <= 0 {
		filtro.RaioMetros = 1500
	}

	filtro.Categoria = strings.TrimSpace(strings.ToLower(filtro.Categoria))

	return filtro
}

func calcularDistanciaMetros(latitudeOrigem float64, longitudeOrigem float64, latitudeDestino float64, longitudeDestino float64) float64 {
	const raioTerraEmMetros = 6371000

	latitudeOrigemRad := grausParaRad(latitudeOrigem)
	latitudeDestinoRad := grausParaRad(latitudeDestino)
	deltaLatitude := grausParaRad(latitudeDestino - latitudeOrigem)
	deltaLongitude := grausParaRad(longitudeDestino - longitudeOrigem)

	a := math.Sin(deltaLatitude/2)*math.Sin(deltaLatitude/2) +
		math.Cos(latitudeOrigemRad)*math.Cos(latitudeDestinoRad)*
			math.Sin(deltaLongitude/2)*math.Sin(deltaLongitude/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return raioTerraEmMetros * c
}

func grausParaRad(graus float64) float64 {
	return graus * math.Pi / 180
}