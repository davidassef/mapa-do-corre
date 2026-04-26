package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/david/mapa-do-corre/backend/internal/models"
)

type CorreRepository struct {
	sqlDB *sql.DB
}

func NewCorreRepository(sqlDB *sql.DB) *CorreRepository {
	return &CorreRepository{sqlDB: sqlDB}
}

func (repository *CorreRepository) BuscarProximos(ctx context.Context, filtro models.FiltroBuscaCorre) ([]models.PrestadorResumo, error) {
	const query = `
		SELECT
			id::text,
			nome,
			categoria,
			descricao,
			whatsapp,
			bairro,
			ST_Y(localizacao::geometry) AS latitude,
			ST_X(localizacao::geometry) AS longitude,
			CAST(ST_Distance(localizacao, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS INTEGER) AS distancia_metros
		FROM prestadores
		WHERE removido_em IS NULL
		  AND ST_DWithin(localizacao, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)
		  AND ($4 = '' OR categoria = $4)
		ORDER BY distancia_metros ASC
	`

	rows, err := repository.sqlDB.QueryContext(ctx, query, filtro.Longitude, filtro.Latitude, filtro.RaioMetros, filtro.Categoria)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar prestadores proximos no postgres: %w", err)
	}
	defer rows.Close()

	prestadores := make([]models.PrestadorResumo, 0)
	for rows.Next() {
		var prestador models.PrestadorResumo
		if err := rows.Scan(
			&prestador.ID,
			&prestador.Nome,
			&prestador.Categoria,
			&prestador.Descricao,
			&prestador.WhatsApp,
			&prestador.Bairro,
			&prestador.Latitude,
			&prestador.Longitude,
			&prestador.DistanciaMetros,
		); err != nil {
			return nil, fmt.Errorf("falha ao ler prestador proximo no postgres: %w", err)
		}

		prestadores = append(prestadores, prestador)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("falha ao iterar prestadores proximos no postgres: %w", err)
	}

	return prestadores, nil
}
