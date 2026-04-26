package service

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/david/mapa-do-corre/backend/internal/models"
)

type CorreStore struct {
	mu                   sync.RWMutex
	prestadores          map[string]prestadorArmazenado
	totalCliques         int
	cliquesPorPrestador  map[string]int
	solicitacoesCadastro map[string]solicitacaoCadastroMemoria
	solicitacoesRemocao  map[string]solicitacaoRemocaoMemoria
}

type prestadorArmazenado struct {
	Resumo     models.PrestadorResumo
	Email      string
	RemovidoEm *time.Time
}

type solicitacaoCadastroMemoria struct {
	Input        models.CadastroPrestadorInput
	CodigoHash   string
	ExpiraEm     time.Time
	ConfirmadoEm *time.Time
	Tentativas   int
}

type solicitacaoRemocaoMemoria struct {
	PrestadorID  string
	Email        string
	CodigoHash   string
	ExpiraEm     time.Time
	ConfirmadoEm *time.Time
	Tentativas   int
}

func NewCorreStore() *CorreStore {
	store := &CorreStore{
		prestadores:          make(map[string]prestadorArmazenado),
		cliquesPorPrestador:  make(map[string]int),
		solicitacoesCadastro: make(map[string]solicitacaoCadastroMemoria),
		solicitacoesRemocao:  make(map[string]solicitacaoRemocaoMemoria),
	}

	store.seedPrestadores()

	return store
}

func (store *CorreStore) ListarPrestadores() []models.PrestadorResumo {
	store.mu.RLock()
	defer store.mu.RUnlock()

	prestadores := make([]models.PrestadorResumo, 0, len(store.prestadores))
	for _, prestador := range store.prestadores {
		if prestador.RemovidoEm != nil {
			continue
		}

		prestadores = append(prestadores, prestador.Resumo)
	}

	sort.Slice(prestadores, func(leftIndex int, rightIndex int) bool {
		return prestadores[leftIndex].Nome < prestadores[rightIndex].Nome
	})

	return prestadores
}

func (store *CorreStore) CriarSolicitacaoCadastro(input models.CadastroPrestadorInput, codigoHash string, expiraEm time.Time) string {
	store.mu.Lock()
	defer store.mu.Unlock()

	solicitacaoID := strconv.FormatInt(time.Now().UnixNano(), 10)
	store.solicitacoesCadastro[solicitacaoID] = solicitacaoCadastroMemoria{
		Input:      input,
		CodigoHash: codigoHash,
		ExpiraEm:   expiraEm,
	}

	return solicitacaoID
}

func (store *CorreStore) LimparSolicitacoesExpiradas() {
	store.mu.Lock()
	defer store.mu.Unlock()

	agora := time.Now()
	for solicitacaoID, solicitacao := range store.solicitacoesCadastro {
		if agora.After(solicitacao.ExpiraEm) || solicitacao.ConfirmadoEm != nil {
			delete(store.solicitacoesCadastro, solicitacaoID)
		}
	}

	for solicitacaoID, solicitacao := range store.solicitacoesRemocao {
		if agora.After(solicitacao.ExpiraEm) || solicitacao.ConfirmadoEm != nil {
			delete(store.solicitacoesRemocao, solicitacaoID)
		}
	}
}

func (store *CorreStore) ConfirmarCadastro(input models.ConfirmacaoCodigoInput, codigoHash string, maxTentativas int) (models.PrestadorResumo, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	solicitacao, exists := store.solicitacoesCadastro[input.SolicitacaoID]
	if !exists {
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoInvalido
	}

	if solicitacao.ConfirmadoEm != nil {
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoInvalido
	}

	if solicitacao.Tentativas >= maxTentativas {
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoExcedeuTentativas
	}

	if time.Now().After(solicitacao.ExpiraEm) {
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoExpirado
	}

	if solicitacao.CodigoHash != codigoHash {
		solicitacao.Tentativas++
		store.solicitacoesCadastro[input.SolicitacaoID] = solicitacao
		return models.PrestadorResumo{}, models.ErrCodigoVerificacaoInvalido
	}

	delete(store.solicitacoesCadastro, input.SolicitacaoID)

	return store.criarPrestadorConfirmado(solicitacao.Input), nil
}

func (store *CorreStore) CriarSolicitacaoRemocao(prestadorID string, email string, codigoHash string, expiraEm time.Time) (string, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()

	prestador, exists := store.prestadores[prestadorID]
	deveEnviar := exists && prestador.RemovidoEm == nil && strings.EqualFold(prestador.Email, email)

	solicitacaoID := strconv.FormatInt(time.Now().UnixNano(), 10)
	solicitacao := solicitacaoRemocaoMemoria{
		Email:      email,
		CodigoHash: codigoHash,
		ExpiraEm:   expiraEm,
	}
	if deveEnviar {
		solicitacao.PrestadorID = prestadorID
	}

	store.solicitacoesRemocao[solicitacaoID] = solicitacao

	return solicitacaoID, deveEnviar
}

func (store *CorreStore) ConfirmarRemocao(input models.ConfirmacaoCodigoInput, codigoHash string, maxTentativas int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	solicitacao, exists := store.solicitacoesRemocao[input.SolicitacaoID]
	if !exists {
		return models.ErrCodigoVerificacaoInvalido
	}

	if solicitacao.ConfirmadoEm != nil {
		return models.ErrCodigoVerificacaoInvalido
	}

	if solicitacao.Tentativas >= maxTentativas {
		return models.ErrCodigoVerificacaoExcedeuTentativas
	}

	if time.Now().After(solicitacao.ExpiraEm) {
		return models.ErrCodigoVerificacaoExpirado
	}

	if solicitacao.CodigoHash != codigoHash {
		solicitacao.Tentativas++
		store.solicitacoesRemocao[input.SolicitacaoID] = solicitacao
		return models.ErrCodigoVerificacaoInvalido
	}

	agora := time.Now()
	delete(store.solicitacoesRemocao, input.SolicitacaoID)

	if solicitacao.PrestadorID == "" {
		return nil
	}

	prestador, exists := store.prestadores[solicitacao.PrestadorID]
	if !exists || prestador.RemovidoEm != nil {
		return nil
	}

	prestador.RemovidoEm = &agora
	store.prestadores[solicitacao.PrestadorID] = prestador

	return nil
}

func (store *CorreStore) RegistrarClique(prestadorID string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	prestador, exists := store.prestadores[prestadorID]
	if !exists || prestador.RemovidoEm != nil {
		return models.ErrPrestadorNaoEncontrado
	}

	store.totalCliques++
	store.cliquesPorPrestador[prestadorID]++

	return nil
}

func (store *CorreStore) ObterResumoImpacto() models.ResumoImpacto {
	store.mu.RLock()
	defer store.mu.RUnlock()

	totaisPorCategoria := make(map[string]models.ConexoesPorCategoria)
	totalPrestadoresAtivos := 0
	totalPrestadoresRemovidos := 0
	for prestadorID, prestador := range store.prestadores {
		totalCategoria := totaisPorCategoria[prestador.Resumo.Categoria]
		totalCategoria.Categoria = prestador.Resumo.Categoria
		if prestador.RemovidoEm == nil {
			totalCategoria.TotalPrestadores++
			totalPrestadoresAtivos++
		} else {
			totalPrestadoresRemovidos++
		}
		totalCategoria.TotalCliques += store.cliquesPorPrestador[prestadorID]
		totaisPorCategoria[prestador.Resumo.Categoria] = totalCategoria
	}

	categorias := make([]models.ConexoesPorCategoria, 0, len(totaisPorCategoria))
	for _, categoria := range totaisPorCategoria {
		categorias = append(categorias, categoria)
	}

	sort.Slice(categorias, func(leftIndex int, rightIndex int) bool {
		return categorias[leftIndex].Categoria < categorias[rightIndex].Categoria
	})

	return models.ResumoImpacto{
		TotalConexoesGeradas:      store.totalCliques,
		TotalPrestadores:          totalPrestadoresAtivos,
		TotalPrestadoresRemovidos: totalPrestadoresRemovidos,
		Categorias:                categorias,
	}
}

func (store *CorreStore) criarPrestadorConfirmado(input models.CadastroPrestadorInput) models.PrestadorResumo {
	prestador := models.PrestadorResumo{
		ID:        strconv.FormatInt(time.Now().UnixNano(), 10),
		Nome:      input.Nome,
		Categoria: input.Categoria,
		Descricao: input.Descricao,
		WhatsApp:  input.WhatsApp,
		Bairro:    input.Bairro,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	}

	store.prestadores[prestador.ID] = prestadorArmazenado{
		Resumo: prestador,
		Email:  input.Email,
	}

	return prestador
}

func (store *CorreStore) seedPrestadores() {
	seedPrestadores := []models.PrestadorResumo{
		{
			ID:        "prestador-001",
			Nome:      "Marmita da Josi",
			Categoria: "alimentacao",
			Descricao: "Marmitas caseiras para almoco em dias uteis.",
			WhatsApp:  "5585988112233",
			Bairro:    "Benfica",
			Latitude:  -3.743269,
			Longitude: -38.536936,
		},
		{
			ID:        "prestador-002",
			Nome:      "Seu Naldo Encanador",
			Categoria: "encanador",
			Descricao: "Atendimento residencial para vazamentos e troca de torneiras.",
			WhatsApp:  "5585988223344",
			Bairro:    "Montese",
			Latitude:  -3.767527,
			Longitude: -38.545087,
		},
		{
			ID:        "prestador-003",
			Nome:      "Dona Liduina Costuras",
			Categoria: "costura",
			Descricao: "Ajustes de roupa e pequenos consertos sob medida.",
			WhatsApp:  "5585988334455",
			Bairro:    "Parquelandia",
			Latitude:  -3.744760,
			Longitude: -38.559191,
		},
		{
			ID:        "prestador-004",
			Nome:      "Bikeboy do Centro",
			Categoria: "entregas",
			Descricao: "Entregas rapidas para pequenos volumes no Centro e arredores.",
			WhatsApp:  "5585988445566",
			Bairro:    "Centro",
			Latitude:  -3.727493,
			Longitude: -38.526670,
		},
		{
			ID:        "prestador-005",
			Nome:      "Rita Faxina Express",
			Categoria: "diarista",
			Descricao: "Faxina residencial por diaria com horario flexivel.",
			WhatsApp:  "5585988556677",
			Bairro:    "Farias Brito",
			Latitude:  -3.734215,
			Longitude: -38.548907,
		},
	}

	for _, prestador := range seedPrestadores {
		store.prestadores[prestador.ID] = prestadorArmazenado{
			Resumo: prestador,
		}
	}
}
