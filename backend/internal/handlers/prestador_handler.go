package handlers

import (
	"errors"
	"strings"

	"github.com/david/mapa-do-corre/backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type prestadorService interface {
	SolicitarCadastro(input models.CadastroPrestadorInput) (models.SolicitacaoCodigoOutput, error)
	ConfirmarCadastro(input models.ConfirmacaoCodigoInput) (models.PrestadorResumo, error)
	SolicitarRemocao(prestadorID string, input models.SolicitacaoRemocaoPrestadorInput) (models.SolicitacaoCodigoOutput, error)
	ConfirmarRemocao(input models.ConfirmacaoCodigoInput) error
	RegistrarClique(prestadorID string, input models.RegistroCliqueInput) error
}

type PrestadorHandler struct {
	service prestadorService
}

func NewPrestadorHandler(service prestadorService) PrestadorHandler {
	return PrestadorHandler{service: service}
}

func (handler PrestadorHandler) Cadastrar(c *fiber.Ctx) error {
	var input models.CadastroPrestadorInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "payload de cadastro invalido",
		})
	}

	solicitacao, err := handler.service.SolicitarCadastro(input)
	if err != nil {
		return responderErroPrestador(c, err)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":     montarMensagemSolicitacaoCadastro(solicitacao.CanalEntrega),
		"solicitacao": solicitacao,
	})
}

func (handler PrestadorHandler) ConfirmarCadastro(c *fiber.Ctx) error {
	var input models.ConfirmacaoCodigoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "payload de confirmacao invalido",
		})
	}

	prestador, err := handler.service.ConfirmarCadastro(input)
	if err != nil {
		return responderErroPrestador(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"prestador": prestador,
	})
}

func (handler PrestadorHandler) RegistrarClique(c *fiber.Ctx) error {
	prestadorID := strings.TrimSpace(c.Params("id"))
	if prestadorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id do prestador e obrigatorio",
		})
	}

	var input models.RegistroCliqueInput
	if c.Body() != nil && len(c.Body()) > 0 {
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "payload de tracking invalido",
			})
		}
	}

	err := handler.service.RegistrarClique(prestadorID, input)
	if err == nil {
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "clique registrado com sucesso",
		})
	}

	if errors.Is(err, models.ErrPrestadorNaoEncontrado) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": "falha ao registrar clique",
	})
}

func (handler PrestadorHandler) SolicitarRemocao(c *fiber.Ctx) error {
	prestadorID := strings.TrimSpace(c.Params("id"))
	if prestadorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id do prestador e obrigatorio",
		})
	}

	var input models.SolicitacaoRemocaoPrestadorInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "payload de remocao invalido",
		})
	}

	solicitacao, err := handler.service.SolicitarRemocao(prestadorID, input)
	if err != nil {
		return responderErroPrestador(c, err)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":     montarMensagemSolicitacaoRemocao(solicitacao.CanalEntrega),
		"solicitacao": solicitacao,
	})
}

func (handler PrestadorHandler) ConfirmarRemocao(c *fiber.Ctx) error {
	var input models.ConfirmacaoCodigoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "payload de confirmacao invalido",
		})
	}

	err := handler.service.ConfirmarRemocao(input)
	if err != nil {
		return responderErroPrestador(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "remocao confirmada com sucesso",
	})
}

func responderErroPrestador(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, models.ErrServicoEmailIndisponivel):
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message": err.Error(),
		})
	case errors.Is(err, models.ErrCodigoVerificacaoExpirado):
		return c.Status(fiber.StatusGone).JSON(fiber.Map{
			"message": err.Error(),
		})
	case errors.Is(err, models.ErrCodigoVerificacaoInvalido), errors.Is(err, models.ErrCodigoVerificacaoExcedeuTentativas):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	case errors.Is(err, models.ErrPrestadorNaoEncontrado):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
}

func montarMensagemSolicitacaoCadastro(canalEntrega string) string {
	if canalEntrega == "debug" {
		return "Codigo gerado em modo local. Use o codigo exibido abaixo para concluir a publicacao."
	}

	return "Enviamos um codigo para o e-mail informado. Confirme-o para publicar o corre."
}

func montarMensagemSolicitacaoRemocao(canalEntrega string) string {
	if canalEntrega == "debug" {
		return "Codigo gerado em modo local. Use o codigo exibido abaixo para concluir a remocao."
	}

	return "Se os dados coincidirem, enviaremos um codigo para o e-mail informado."
}
