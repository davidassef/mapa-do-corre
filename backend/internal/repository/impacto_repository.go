package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/david/mapa-do-corre/backend/internal/models"
)

type ImpactoRepository struct {
	sqlDB *sql.DB
}

func NewImpactoRepository(sqlDB *sql.DB) *ImpactoRepository {
	return &ImpactoRepository{sqlDB: sqlDB}
}

func (repository *ImpactoRepository) ObterResumo(ctx context.Context) (models.ResumoImpacto, error) {
	const resumoQuery = `
		SELECT
			(SELECT COUNT(*) FROM prestadores WHERE removido_em IS NULL) AS total_prestadores,
			(SELECT COUNT(*) FROM prestadores WHERE removido_em IS NOT NULL) AS total_prestadores_removidos,
			(SELECT COUNT(*) FROM logs_cliques) AS total_conexoes
	`

	var resumo models.ResumoImpacto
	if err := repository.sqlDB.QueryRowContext(ctx, resumoQuery).Scan(
		&resumo.TotalPrestadores,
		&resumo.TotalPrestadoresRemovidos,
		&resumo.TotalConexoesGeradas,
	); err != nil {
		return models.ResumoImpacto{}, fmt.Errorf("falha ao carregar resumo geral no postgres: %w", err)
	}

	const categoriasQuery = `
		WITH prestadores_ativos AS (
			SELECT categoria, COUNT(*) AS total_prestadores
			FROM prestadores
			WHERE removido_em IS NULL
			GROUP BY categoria
		),
		cliques_historicos AS (
			SELECT p.categoria, COUNT(l.id) AS total_cliques
			FROM prestadores p
			LEFT JOIN logs_cliques l ON l.prestador_id = p.id
			GROUP BY p.categoria
		)
		SELECT
			COALESCE(a.categoria, h.categoria) AS categoria,
			COALESCE(a.total_prestadores, 0) AS total_prestadores,
			COALESCE(h.total_cliques, 0) AS total_cliques
		FROM prestadores_ativos a
		FULL OUTER JOIN cliques_historicos h ON h.categoria = a.categoria
		ORDER BY categoria ASC
	`

	rows, err := repository.sqlDB.QueryContext(ctx, categoriasQuery)
	if err != nil {
		return models.ResumoImpacto{}, fmt.Errorf("falha ao carregar resumo por categoria no postgres: %w", err)
	}
	defer rows.Close()

	resumo.Categorias = make([]models.ConexoesPorCategoria, 0)
	for rows.Next() {
		var categoria models.ConexoesPorCategoria
		if err := rows.Scan(
			&categoria.Categoria,
			&categoria.TotalPrestadores,
			&categoria.TotalCliques,
		); err != nil {
			return models.ResumoImpacto{}, fmt.Errorf("falha ao ler categoria de impacto no postgres: %w", err)
		}

		resumo.Categorias = append(resumo.Categorias, categoria)
	}

	if err := rows.Err(); err != nil {
		return models.ResumoImpacto{}, fmt.Errorf("falha ao iterar categorias de impacto no postgres: %w", err)
	}

	return resumo, nil
}
