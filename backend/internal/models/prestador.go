package models

import "time"

type FiltroBuscaCorre struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	RaioMetros int     `json:"raioMetros"`
	Categoria  string  `json:"categoria,omitempty"`
}

type PrestadorResumo struct {
	ID              string  `json:"id"`
	Nome            string  `json:"nome"`
	Categoria       string  `json:"categoria"`
	Descricao       string  `json:"descricao"`
	WhatsApp        string  `json:"whatsApp"`
	Bairro          string  `json:"bairro"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	DistanciaMetros int     `json:"distanciaMetros"`
}

type CadastroPrestadorInput struct {
	Nome      string  `json:"nome"`
	Categoria string  `json:"categoria"`
	Descricao string  `json:"descricao"`
	WhatsApp  string  `json:"whatsApp"`
	Bairro    string  `json:"bairro"`
	Email     string  `json:"email"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ConfirmacaoCodigoInput struct {
	SolicitacaoID string `json:"solicitacaoId"`
	Codigo        string `json:"codigo"`
}

type SolicitacaoRemocaoPrestadorInput struct {
	Email string `json:"email"`
}

type SolicitacaoCodigoOutput struct {
	SolicitacaoID string    `json:"solicitacaoId"`
	ExpiraEm      time.Time `json:"expiraEm"`
	CanalEntrega  string    `json:"canalEntrega"`
	DebugCodigo   string    `json:"debugCodigo,omitempty"`
}

type RegistroCliqueInput struct {
	Origem string `json:"origem"`
}
