package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/david/mapa-do-corre/backend/internal/models"
)

type prestadorRepository interface {
	LimparSolicitacoesExpiradas(ctx context.Context) error
	CriarSolicitacaoCadastro(ctx context.Context, input models.CadastroPrestadorInput, codigoHash string, expiraEm time.Time) (string, error)
	ConfirmarCadastro(ctx context.Context, input models.ConfirmacaoCodigoInput, codigoHash string, maxTentativas int) (models.PrestadorResumo, error)
	CriarSolicitacaoRemocao(ctx context.Context, prestadorID string, email string, codigoHash string, expiraEm time.Time) (string, bool, error)
	ConfirmarRemocao(ctx context.Context, input models.ConfirmacaoCodigoInput, codigoHash string, maxTentativas int) error
	RegistrarClique(ctx context.Context, prestadorID string, origem string) error
}

type PrestadorService struct {
	store                  *CorreStore
	repository             prestadorRepository
	emailSender            EmailSender
	exibeCodigoDebug       bool
	emailCodeTTL           time.Duration
	emailCodeMaxTentativas int
}

func NewPrestadorService(
	store *CorreStore,
	repository prestadorRepository,
	emailSender EmailSender,
	exibeCodigoDebug bool,
	emailCodeTTL time.Duration,
	emailCodeMaxTentativas int,
) *PrestadorService {
	if emailCodeTTL <= 0 {
		emailCodeTTL = 10 * time.Minute
	}

	if emailCodeMaxTentativas <= 0 {
		emailCodeMaxTentativas = 5
	}

	return &PrestadorService{
		store:                  store,
		repository:             repository,
		emailSender:            emailSender,
		exibeCodigoDebug:       exibeCodigoDebug,
		emailCodeTTL:           emailCodeTTL,
		emailCodeMaxTentativas: emailCodeMaxTentativas,
	}
}

func (service *PrestadorService) SolicitarCadastro(input models.CadastroPrestadorInput) (models.SolicitacaoCodigoOutput, error) {
	service.limparSolicitacoesExpiradas()

	if err := service.validarServicoEmail(); err != nil {
		return models.SolicitacaoCodigoOutput{}, err
	}

	input, err := normalizarCadastroInput(input)
	if err != nil {
		return models.SolicitacaoCodigoOutput{}, err
	}

	codigo, err := gerarCodigoNumerico(6)
	if err != nil {
		return models.SolicitacaoCodigoOutput{}, fmt.Errorf("falha ao gerar codigo de verificacao: %w", err)
	}

	grupoExpiracao := time.Now().Add(service.emailCodeTTL)
	codigoHash := gerarHashCodigo(codigo)

	var solicitacaoID string
	if service.repository != nil {
		commandCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		solicitacaoID, err = service.repository.CriarSolicitacaoCadastro(commandCtx, input, codigoHash, grupoExpiracao)
		if err != nil {
			return models.SolicitacaoCodigoOutput{}, err
		}
	} else {
		solicitacaoID = service.store.CriarSolicitacaoCadastro(input, codigoHash, grupoExpiracao)
	}

	if err := service.enviarCodigoPorEmail(input.Email, "Confirme seu corre no Mapa do Corre", montarMensagemCodigoCadastro(codigo, grupoExpiracao)); err != nil {
		return models.SolicitacaoCodigoOutput{}, err
	}

	return service.montarSolicitacaoCodigo(solicitacaoID, grupoExpiracao, codigo, true), nil
}

func (service *PrestadorService) ConfirmarCadastro(input models.ConfirmacaoCodigoInput) (models.PrestadorResumo, error) {
	service.limparSolicitacoesExpiradas()

	input, codigoHash, err := service.normalizarConfirmacaoCodigo(input)
	if err != nil {
		return models.PrestadorResumo{}, err
	}

	if service.repository != nil {
		commandCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return service.repository.ConfirmarCadastro(commandCtx, input, codigoHash, service.emailCodeMaxTentativas)
	}

	return service.store.ConfirmarCadastro(input, codigoHash, service.emailCodeMaxTentativas)
}

func (service *PrestadorService) RegistrarClique(prestadorID string, input models.RegistroCliqueInput) error {
	origem := strings.TrimSpace(input.Origem)
	if origem == "" {
		origem = "whatsapp"
	}

	if service.repository != nil {
		commandCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return service.repository.RegistrarClique(commandCtx, prestadorID, origem)
	}

	return service.store.RegistrarClique(prestadorID)
}

func (service *PrestadorService) SolicitarRemocao(prestadorID string, input models.SolicitacaoRemocaoPrestadorInput) (models.SolicitacaoCodigoOutput, error) {
	service.limparSolicitacoesExpiradas()

	if err := service.validarServicoEmail(); err != nil {
		return models.SolicitacaoCodigoOutput{}, err
	}

	prestadorID = strings.TrimSpace(prestadorID)
	if prestadorID == "" {
		return models.SolicitacaoCodigoOutput{}, fmt.Errorf("id do prestador e obrigatorio")
	}

	email, err := normalizarEmail(input.Email)
	if err != nil {
		return models.SolicitacaoCodigoOutput{}, err
	}

	codigo, err := gerarCodigoNumerico(6)
	if err != nil {
		return models.SolicitacaoCodigoOutput{}, fmt.Errorf("falha ao gerar codigo de verificacao: %w", err)
	}

	grupoExpiracao := time.Now().Add(service.emailCodeTTL)
	codigoHash := gerarHashCodigo(codigo)

	var (
		solicitacaoID string
		deveEnviar    bool
	)

	if service.repository != nil {
		commandCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		solicitacaoID, deveEnviar, err = service.repository.CriarSolicitacaoRemocao(commandCtx, prestadorID, email, codigoHash, grupoExpiracao)
		if err != nil {
			return models.SolicitacaoCodigoOutput{}, err
		}
	} else {
		solicitacaoID, deveEnviar = service.store.CriarSolicitacaoRemocao(prestadorID, email, codigoHash, grupoExpiracao)
	}

	if deveEnviar {
		if err := service.enviarCodigoPorEmail(email, "Confirme a remocao do seu corre", montarMensagemCodigoRemocao(codigo, grupoExpiracao)); err != nil {
			return models.SolicitacaoCodigoOutput{}, err
		}
	}

	return service.montarSolicitacaoCodigo(solicitacaoID, grupoExpiracao, codigo, deveEnviar), nil
}

func (service *PrestadorService) ConfirmarRemocao(input models.ConfirmacaoCodigoInput) error {
	service.limparSolicitacoesExpiradas()

	input, codigoHash, err := service.normalizarConfirmacaoCodigo(input)
	if err != nil {
		return err
	}

	if service.repository != nil {
		commandCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return service.repository.ConfirmarRemocao(commandCtx, input, codigoHash, service.emailCodeMaxTentativas)
	}

	return service.store.ConfirmarRemocao(input, codigoHash, service.emailCodeMaxTentativas)
}

func (service *PrestadorService) limparSolicitacoesExpiradas() {
	if service.repository != nil {
		commandCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := service.repository.LimparSolicitacoesExpiradas(commandCtx); err != nil {
			log.Printf("[WARN] Falha ao limpar solicitacoes expiradas no postgres: %v", err)
		}

		return
	}

	service.store.LimparSolicitacoesExpiradas()
}

func (service *PrestadorService) validarServicoEmail() error {
	if service.emailSender != nil || service.exibeCodigoDebug {
		return nil
	}

	return models.ErrServicoEmailIndisponivel
}

func (service *PrestadorService) enviarCodigoPorEmail(email string, assunto string, mensagem string) error {
	if service.exibeCodigoDebug {
		log.Printf("[INFO] Codigo de verificacao em modo local para %s: %s", email, extrairCodigoDaMensagem(mensagem))
		return nil
	}

	if service.emailSender == nil {
		return models.ErrServicoEmailIndisponivel
	}

	commandCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return service.emailSender.Enviar(commandCtx, EmailMessage{
		Destinatario: email,
		Assunto:      assunto,
		Corpo:        mensagem,
	})
}

func (service *PrestadorService) montarSolicitacaoCodigo(solicitacaoID string, expiraEm time.Time, codigo string, exibirCodigo bool) models.SolicitacaoCodigoOutput {
	resultado := models.SolicitacaoCodigoOutput{
		SolicitacaoID: solicitacaoID,
		ExpiraEm:      expiraEm,
		CanalEntrega:  "email",
	}

	if service.exibeCodigoDebug && exibirCodigo {
		resultado.CanalEntrega = "debug"
		resultado.DebugCodigo = codigo
	}

	return resultado
}

func (service *PrestadorService) normalizarConfirmacaoCodigo(input models.ConfirmacaoCodigoInput) (models.ConfirmacaoCodigoInput, string, error) {
	input.SolicitacaoID = strings.TrimSpace(input.SolicitacaoID)
	if input.SolicitacaoID == "" {
		return models.ConfirmacaoCodigoInput{}, "", fmt.Errorf("solicitacaoId e obrigatorio")
	}

	builder := strings.Builder{}
	for _, caractere := range input.Codigo {
		if !unicode.IsDigit(caractere) {
			continue
		}

		builder.WriteRune(caractere)
	}

	input.Codigo = builder.String()
	if len(input.Codigo) != 6 {
		return models.ConfirmacaoCodigoInput{}, "", fmt.Errorf("codigo deve ter 6 digitos")
	}

	return input, gerarHashCodigo(input.Codigo), nil
}

func normalizarCadastroInput(input models.CadastroPrestadorInput) (models.CadastroPrestadorInput, error) {
	if strings.TrimSpace(input.Nome) == "" {
		return models.CadastroPrestadorInput{}, fmt.Errorf("nome e obrigatorio")
	}

	if strings.TrimSpace(input.Categoria) == "" {
		return models.CadastroPrestadorInput{}, fmt.Errorf("categoria e obrigatoria")
	}

	if input.Latitude == 0 || input.Longitude == 0 {
		return models.CadastroPrestadorInput{}, fmt.Errorf("latitude e longitude sao obrigatorias")
	}

	input.WhatsApp = normalizarWhatsApp(input.WhatsApp)
	if input.WhatsApp == "" {
		return models.CadastroPrestadorInput{}, fmt.Errorf("whatsApp e obrigatorio")
	}

	email, err := normalizarEmail(input.Email)
	if err != nil {
		return models.CadastroPrestadorInput{}, err
	}

	input.Nome = strings.TrimSpace(input.Nome)
	input.Categoria = strings.TrimSpace(strings.ToLower(input.Categoria))
	input.Descricao = strings.TrimSpace(input.Descricao)
	input.Bairro = strings.TrimSpace(input.Bairro)
	input.Email = email

	return input, nil
}

func normalizarEmail(rawEmail string) (string, error) {
	endereco, err := mail.ParseAddress(strings.TrimSpace(rawEmail))
	if err != nil {
		return "", fmt.Errorf("email invalido")
	}

	return strings.ToLower(strings.TrimSpace(endereco.Address)), nil
}

func gerarCodigoNumerico(tamanho int) (string, error) {
	if tamanho <= 0 {
		return "", fmt.Errorf("tamanho de codigo invalido")
	}

	bytes := make([]byte, tamanho)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	builder := strings.Builder{}
	for _, byteAtual := range bytes {
		builder.WriteByte('0' + (byteAtual % 10))
	}

	return builder.String(), nil
}

func gerarHashCodigo(codigo string) string {
	sum := sha256.Sum256([]byte(codigo))
	return hex.EncodeToString(sum[:])
}

func montarMensagemCodigoCadastro(codigo string, expiraEm time.Time) string {
	return fmt.Sprintf(
		"Seu codigo para publicar o corre no Mapa do Corre e: %s\n\nEle expira em %s.",
		codigo,
		expiraEm.Format("02/01/2006 15:04"),
	)
}

func montarMensagemCodigoRemocao(codigo string, expiraEm time.Time) string {
	return fmt.Sprintf(
		"Seu codigo para confirmar a remocao do corre no Mapa do Corre e: %s\n\nEle expira em %s.",
		codigo,
		expiraEm.Format("02/01/2006 15:04"),
	)
}

func extrairCodigoDaMensagem(mensagem string) string {
	partes := strings.Split(mensagem, ":")
	if len(partes) < 2 {
		return ""
	}

	resto := strings.TrimSpace(partes[1])
	if indiceQuebra := strings.Index(resto, "\n"); indiceQuebra >= 0 {
		return strings.TrimSpace(resto[:indiceQuebra])
	}

	return resto
}

func normalizarWhatsApp(rawPhone string) string {
	builder := strings.Builder{}
	for _, character := range rawPhone {
		if !unicode.IsDigit(character) {
			continue
		}

		builder.WriteRune(character)
	}

	whatsApp := builder.String()
	if whatsApp == "" {
		return ""
	}

	if len(whatsApp) == 11 {
		return "55" + whatsApp
	}

	return whatsApp
}
