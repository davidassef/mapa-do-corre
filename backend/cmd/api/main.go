package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/david/mapa-do-corre/backend/internal/config"
	"github.com/david/mapa-do-corre/backend/internal/database"
	"github.com/david/mapa-do-corre/backend/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("[ERROR] falha ao carregar configuracao: %v", err)
	}

	databaseCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	databaseConnection, err := database.NewConnection(databaseCtx, appConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("[ERROR] falha ao conectar no banco: %v", err)
	}

	if databaseConnection != nil {
		defer databaseConnection.Close()
	}

	app := server.NewApp(appConfig, databaseConnection)

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("[ERROR] falha no shutdown da API: %v", err)
		}
	}()

	log.Printf("[INFO] API ouvindo na porta %s", appConfig.Port)

	if err := app.Listen(":" + appConfig.Port); err != nil {
		log.Fatalf("[ERROR] falha ao subir a API: %v", err)
	}
}