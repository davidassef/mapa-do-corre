package models

import "errors"

var ErrCodigoVerificacaoInvalido = errors.New("codigo de verificacao invalido")
var ErrCodigoVerificacaoExpirado = errors.New("codigo de verificacao expirado")
var ErrCodigoVerificacaoExcedeuTentativas = errors.New("limite de tentativas do codigo excedido")
var ErrServicoEmailIndisponivel = errors.New("verificacao por email indisponivel no momento")
