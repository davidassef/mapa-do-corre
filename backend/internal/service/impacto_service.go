package service

import (
	"context"
	"time"

	"github.com/david/mapa-do-corre/backend/internal/models"
)

type impactoReader interface {
	ObterResumo(ctx context.Context) (models.ResumoImpacto, error)
}

type ImpactoService struct {
	store      *CorreStore
	repository impactoReader
}

func NewImpactoService(store *CorreStore, repository impactoReader) *ImpactoService {
	return &ImpactoService{store: store, repository: repository}
}

func (service *ImpactoService) ObterResumo() (models.ResumoImpacto, error) {
	if service.repository != nil {
		queryCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return service.repository.ObterResumo(queryCtx)
	}

	return service.store.ObterResumoImpacto(), nil
}