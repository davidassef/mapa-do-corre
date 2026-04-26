package models

import "errors"

var ErrPrestadorNaoEncontrado = errors.New("prestador nao encontrado")

type ConexoesPorCategoria struct {
	Categoria        string `json:"categoria"`
	TotalPrestadores int    `json:"totalPrestadores"`
	TotalCliques     int    `json:"totalCliques"`
}

type ResumoImpacto struct {
	TotalConexoesGeradas      int                    `json:"totalConexoesGeradas"`
	TotalPrestadores          int                    `json:"totalPrestadores"`
	TotalPrestadoresRemovidos int                    `json:"totalPrestadoresRemovidos"`
	Categorias                []ConexoesPorCategoria `json:"categorias"`
}
