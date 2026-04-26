package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Connection struct {
	sqlDB *sql.DB
}

func NewConnection(ctx context.Context, databaseURL string) (*Connection, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, nil
	}

	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexao com postgres: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("falha ao validar conexao com postgres: %w", err)
	}

	return &Connection{sqlDB: sqlDB}, nil
}

func (connection *Connection) Ping(ctx context.Context) error {
	if connection == nil || connection.sqlDB == nil {
		return nil
	}

	return connection.sqlDB.PingContext(ctx)
}

func (connection *Connection) DB() *sql.DB {
	if connection == nil {
		return nil
	}

	return connection.sqlDB
}

func (connection *Connection) Close() error {
	if connection == nil || connection.sqlDB == nil {
		return nil
	}

	return connection.sqlDB.Close()
}