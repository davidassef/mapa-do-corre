package server

import (
	"time"

	"github.com/david/mapa-do-corre/backend/internal/config"
	"github.com/david/mapa-do-corre/backend/internal/database"
	"github.com/david/mapa-do-corre/backend/internal/handlers"
	"github.com/david/mapa-do-corre/backend/internal/middleware"
	"github.com/david/mapa-do-corre/backend/internal/repository"
	"github.com/david/mapa-do-corre/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	requestid "github.com/gofiber/fiber/v2/middleware/requestid"
)

func NewApp(appConfig config.Config, databaseConnection *database.Connection) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Mapa do Corre API",
	})

	store := service.NewCorreStore()
	var correRepository *repository.CorreRepository
	var prestadorRepository *repository.PrestadorRepository
	var impactoRepository *repository.ImpactoRepository
	var emailSender service.EmailSender

	if appConfig.PersistenceMode == "postgres" && databaseConnection != nil && databaseConnection.DB() != nil {
		sqlDB := databaseConnection.DB()
		correRepository = repository.NewCorreRepository(sqlDB)
		prestadorRepository = repository.NewPrestadorRepository(sqlDB)
		impactoRepository = repository.NewImpactoRepository(sqlDB)
	}

	if appConfig.HasSMTPConfiguration() {
		emailSender = service.NewSMTPEmailSender(
			appConfig.SMTPHost,
			appConfig.SMTPPort,
			appConfig.SMTPUsername,
			appConfig.SMTPPassword,
			appConfig.SMTPFromEmail,
			appConfig.SMTPFromName,
		)
	}

	correService := service.NewCorreService(store, correRepository)
	prestadorService := service.NewPrestadorService(
		store,
		prestadorRepository,
		emailSender,
		!appConfig.HasSMTPConfiguration() && appConfig.AppEnv != "production",
		time.Duration(appConfig.EmailCodeTTLMinutes)*time.Minute,
		appConfig.EmailCodeMaxAttempts,
	)
	impactoService := service.NewImpactoService(store, impactoRepository)

	healthHandler := handlers.NewHealthHandler(databaseConnection)
	correHandler := handlers.NewCorreHandler(correService)
	prestadorHandler := handlers.NewPrestadorHandler(prestadorService)
	impactoHandler := handlers.NewImpactoHandler(impactoService)

	app.Use(requestid.New())
	app.Use(middleware.NewCors(appConfig.AppEnv, appConfig.FrontendOrigin))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "Mapa do Corre API",
			"version": "0.1.0",
		})
	})

	app.Get("/health", healthHandler.Handle)
	app.Get("/corres", correHandler.Buscar)
	app.Get("/impacto/resumo", impactoHandler.Resumo)

	publicWriteRoutes := app.Group("/", middleware.NewPublicWriteLimiter(appConfig.WriteRateLimitMax, appConfig.WriteRateLimitWindowSeconds))
	publicWriteRoutes.Post("/prestadores", prestadorHandler.Cadastrar)
	publicWriteRoutes.Post("/prestadores/confirmacoes", prestadorHandler.ConfirmarCadastro)
	publicWriteRoutes.Post("/prestadores/:id/cliques", prestadorHandler.RegistrarClique)
	publicWriteRoutes.Post("/prestadores/:id/remocao", prestadorHandler.SolicitarRemocao)
	publicWriteRoutes.Post("/prestadores/remocoes/confirmacoes", prestadorHandler.ConfirmarRemocao)

	return app
}
