package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/david/mapa-do-corre/backend/internal/models"
	"github.com/google/uuid"
)

type PrestadorRepository struct {
	sqlDB *sql.DB
}

func NewPrestadorRepository(sqlDB *sql.DB) *PrestadorRepository {
	return &PrestadorRepository{sqlDB: sqlDB}
}

func (repository *PrestadorRepository) LimparSolicitacoesExpiradas(ctx context.Context) error {
	const query = `
		DELETE FROM solicitacoes_cadastro_prestador
		WHERE expira_em < NOW()
		   OR confirmado_em IS NOT NULL;

		DELETE FROM solicitacoes_remocao_prestador
		WHERE expira_em < NOW()
		   OR confirmado_em IS NOT NULL;
	`

	if _, err := repository.sqlDB.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("falha ao limpar solicitacoes expiradas no postgres: %w", err)
	}

	return nil
}

func (repository *PrestadorRepository) CriarSolicitacaoCadastro(ctx context.Context, input models.CadastroPrestadorInput, codigoHash string, expiraEm time.Time) (string, error) {
	const query = `
		INSERT INTO solicitacoes_cadastro_prestador (
			id,
			nome,
			categoria,
			descricao,
			whatsapp,
			bairro,
			email,
			latitude,
			longitude,
			codigo_hash,
			expira_em
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11
		)
	`

	solicitacaoID := uuid.NewString()
	if _, err := repository.sqlDB.ExecContext(
		ctx,
		query,
		solicitacaoID,
		input.Nome,
		input.Categoria,
		input.Descricao,
		input.WhatsApp,
		input.Bairro,
		input.Email,
		input.Latitude,
		input.Longitude,
		codigoHash,
		expiraEm,
	); err != nil {
		return "", fmt.Errorf("falha ao criar solicitacao de cadastro no postgres: %w", err)
	}

	return solicitacaoID, nil
}

func (repository *PrestadorRepository) ConfirmarCadastro(ctx context.Context, input models.ConfirmacaoCodigoInput, codigoHash string, maxTentativas int) (models.PrestadorResumo, error) {
	tx, err := repository.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return models.PrestadorResumo{}, fmt.Errorf("falha ao abrir transacao de cadastro: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	const selectQuery = `
		SELECT
			nome,
			categoria,
			descricao,
			whatsapp,
			bairro,
			email,
			latitude,
			longitude,
			codigo_hash,
			expira_em,
			confirmado_em,
			tentativas
		FROM solicitacoes_cadastro_prestador
		WHERE id = $1
		FOR UPDATE
	`

	var solicitacao struct {
		Nome         string
		Categoria    string
		Descricao    string
		WhatsApp     string
		Bairro       string
		Email        string
		Latitude     float64
		Longitude    float64
		CodigoHash   string
		ExpiraEm     time.Time
		ConfirmadoEm sql.NullTime
		Tentativas   int
	}

	err = tx.QueryRowContext(ctx, selectQuery, input.SolicitacaoID).Scan(
		&solicitacao.Nome,
		&solicitacao.Categoria,
		&solicitacao.Descricao,
		&solicitacao.WhatsApp,
		&solicitacao.Bairro,
		&solicitacao.Email,
		&solicitacao.Latitude,
		&solicitacao.Longitude,
		&solicitacao.CodigoHash,
		&solicitacao.ExpiraEm,
		&solicitacao.ConfirmadoEm,
		&solicitacao.Tentativas,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.PrestadorResumo{}, models.ErrCodigoVerificacaoInvalido
		}

		return models.PrestadorResumo{}, fmt.Errorf("falha ao carregar solicitacao de cadastro: %w", err)
	}

	if solicitacao.ConfirmadoEm.Valid {
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoInvalido
	}

	if solicitacao.Tentativas >= maxTentativas {
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoExcedeuTentativas
	}

	if time.Now().After(solicitacao.ExpiraEm) {
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoExpirado
	}

	if solicitacao.CodigoHash != codigoHash {
		if _, err := tx.ExecContext(ctx, `
			UPDATE solicitacoes_cadastro_prestador
			SET tentativas = tentativas + 1, updated_at = NOW()
			WHERE id = $1
		`, input.SolicitacaoID); err != nil {
			return models.PrestadorResumo{}, fmt.Errorf("falha ao registrar tentativa de cadastro: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return models.PrestadorResumo{}, fmt.Errorf("falha ao confirmar tentativa invalida de cadastro: %w", err)
		}
		committed = true

		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoInvalido
	}

	const insertPrestadorQuery = `
		INSERT INTO prestadores (
			id,
			nome,
			categoria,
			descricao,
			whatsapp,
			bairro,
			email_responsavel,
			latitude,
			longitude
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9
		)
		RETURNING
			id::text,
			nome,
			categoria,
			descricao,
			whatsapp,
			bairro,
			latitude,
			longitude
	`

	prestadorID := uuid.NewString()
	var prestador models.PrestadorResumo
	if err := tx.QueryRowContext(
		ctx,
		insertPrestadorQuery,
		prestadorID,
		solicitacao.Nome,
		solicitacao.Categoria,
		solicitacao.Descricao,
		solicitacao.WhatsApp,
		solicitacao.Bairro,
		solicitacao.Email,
		solicitacao.Latitude,
		solicitacao.Longitude,
	).Scan(
		&prestador.ID,
		&prestador.Nome,
		&prestador.Categoria,
		&prestador.Descricao,
		&prestador.WhatsApp,
		&prestador.Bairro,
		&prestador.Latitude,
		&prestador.Longitude,
	); err != nil {
		return models.PrestadorResumo{}, fmt.Errorf("falha ao criar prestador confirmado no postgres: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM solicitacoes_cadastro_prestador
		WHERE id = $1
	`, input.SolicitacaoID); err != nil {
		return models.PrestadorResumo{}, fmt.Errorf("falha ao limpar solicitacao de cadastro confirmada: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.PrestadorResumo{}, fmt.Errorf("falha ao confirmar cadastro no postgres: %w", err)
	}
	committed = true

	return prestador, nil
}

func (repository *PrestadorRepository) CriarSolicitacaoRemocao(ctx context.Context, prestadorID string, email string, codigoHash string, expiraEm time.Time) (string, bool, error) {
	const lookupQuery = `
		SELECT id::text
		FROM prestadores
		WHERE id = $1
		  AND removido_em IS NULL
		  AND LOWER(email_responsavel) = LOWER($2)
	`

	var prestadorEncontradoID string
	err := repository.sqlDB.QueryRowContext(ctx, lookupQuery, prestadorID, email).Scan(&prestadorEncontradoID)
	deveEnviarCodigo := err == nil
	if err != nil && err != sql.ErrNoRows {
		return "", false, fmt.Errorf("falha ao validar prestador para remocao: %w", err)
	}

	const insertQuery = `
		INSERT INTO solicitacoes_remocao_prestador (
			id,
			prestador_id,
			email,
			codigo_hash,
			expira_em
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	solicitacaoID := uuid.NewString()
	var prestadorIDParam any
	if deveEnviarCodigo {
		prestadorIDParam = prestadorEncontradoID
	}

	if _, err := repository.sqlDB.ExecContext(ctx, insertQuery, solicitacaoID, prestadorIDParam, email, codigoHash, expiraEm); err != nil {
		return "", false, fmt.Errorf("falha ao criar solicitacao de remocao no postgres: %w", err)
	}

	return solicitacaoID, deveEnviarCodigo, nil
}

func (repository *PrestadorRepository) ConfirmarRemocao(ctx context.Context, input models.ConfirmacaoCodigoInput, codigoHash string, maxTentativas int) error {
	tx, err := repository.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("falha ao abrir transacao de remocao: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	const selectQuery = `
		SELECT prestador_id::text, codigo_hash, expira_em, confirmado_em, tentativas
		FROM solicitacoes_remocao_prestador
		WHERE id = $1
		FOR UPDATE
	`

	var prestadorID sql.NullString
	var codigoHashArmazenado string
	var expiraEm time.Time
	var confirmadoEm sql.NullTime
	var tentativas int

	err = tx.QueryRowContext(ctx, selectQuery, input.SolicitacaoID).Scan(
		&prestadorID,
		&codigoHashArmazenado,
		&expiraEm,
		&confirmadoEm,
		&tentativas,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ErrCodigoVerificacaoInvalido
		}

		return fmt.Errorf("falha ao carregar solicitacao de remocao: %w", err)
	}

	if confirmadoEm.Valid {
		return models.ErrCodigoVerificacaoInvalido
	}

	if tentativas >= maxTentativas {
		return models.ErrCodigoVerificacaoExcedeuTentativas
	}

	if time.Now().After(expiraEm) {
		return models.ErrCodigoVerificacaoExpirado
	}

	if codigoHashArmazenado != codigoHash {
		if _, err := tx.ExecContext(ctx, `
			UPDATE solicitacoes_remocao_prestador
			SET tentativas = tentativas + 1, updated_at = NOW()
			WHERE id = $1
		`, input.SolicitacaoID); err != nil {
			return fmt.Errorf("falha ao registrar tentativa de remocao: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("falha ao confirmar tentativa invalida de remocao: %w", err)
		}
		committed = true

		return models.ErrCodigoVerificacaoInvalido
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM solicitacoes_remocao_prestador
		WHERE id = $1
	`, input.SolicitacaoID); err != nil {
		return fmt.Errorf("falha ao limpar solicitacao de remocao confirmada: %w", err)
	}

	if prestadorID.Valid {
		if _, err := tx.ExecContext(ctx, `
			UPDATE prestadores
			SET removido_em = COALESCE(removido_em, NOW()), updated_at = NOW()
			WHERE id = $1
		`, prestadorID.String); err != nil {
			return fmt.Errorf("falha ao remover prestador no postgres: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("falha ao confirmar remocao no postgres: %w", err)
	}
	committed = true

	return nil
}

func (repository *PrestadorRepository) RegistrarClique(ctx context.Context, prestadorID string, origem string) error {
	const query = `
		INSERT INTO logs_cliques (id, prestador_id, origem)
		SELECT $1, id, $3
		FROM prestadores
		WHERE id = $2
		  AND removido_em IS NULL
	`

	result, err := repository.sqlDB.ExecContext(ctx, query, uuid.NewString(), prestadorID, origem)
	if err != nil {
		return fmt.Errorf("falha ao registrar clique no postgres: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao ler resultado do tracking no postgres: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrPrestadorNaoEncontrado
	}

	return nil
}
