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
		WITH candidatos AS (
			SELECT
				id::text,
				nome,
				categoria,
				descricao,
				whatsapp,
				bairro,
				latitude,
				longitude,
				POWER(SIN(RADIANS(latitude - $2) / 2), 2) +
					COS(RADIANS($2)) * COS(RADIANS(latitude)) * POWER(SIN(RADIANS(longitude - $1) / 2), 2) AS haversine_a
			FROM prestadores
			WHERE removido_em IS NULL
			  AND ($4 = '' OR categoria = $4)
			  AND latitude BETWEEN $2 - ($3 / 111320.0) AND $2 + ($3 / 111320.0)
			  AND longitude BETWEEN $1 - ($3 / (111320.0 * GREATEST(COS(RADIANS($2)), 0.01)))
				  AND $1 + ($3 / (111320.0 * GREATEST(COS(RADIANS($2)), 0.01)))
		),
		distancias AS (
			SELECT
				id,
				nome,
				categoria,
				descricao,
				whatsapp,
				bairro,
				latitude,
				longitude,
				CAST(6371000 * 2 * ASIN(SQRT(LEAST(1, haversine_a))) AS INTEGER) AS distancia_metros
			FROM candidatos
		)
		SELECT
			id,
			nome,
			categoria,
			descricao,
			whatsapp,
			bairro,
			latitude,
			longitude,
			distancia_metros
		FROM distancias
		WHERE distancia_metros <= $3
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
