package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                      string
	Port                        string
	DatabaseURL                 string
	FrontendOrigin              string
	PersistenceMode             string
	SMTPHost                    string
	SMTPPort                    string
	SMTPUsername                string
	SMTPPassword                string
	SMTPFromEmail               string
	SMTPFromName                string
	EmailCodeTTLMinutes         int
	EmailCodeMaxAttempts        int
	WriteRateLimitMax           int
	WriteRateLimitWindowSeconds int
}

func Load() (Config, error) {
	loadEnvFiles()

	appConfig := Config{
		AppEnv:                      readEnv("APP_ENV", "development"),
		Port:                        readEnv("PORT", "8080"),
		DatabaseURL:                 strings.TrimSpace(os.Getenv("DATABASE_URL")),
		FrontendOrigin:              readEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
		PersistenceMode:             readEnv("PERSISTENCE_MODE", "memory"),
		SMTPHost:                    readEnv("SMTP_HOST", ""),
		SMTPPort:                    readEnv("SMTP_PORT", ""),
		SMTPUsername:                readEnvAny([]string{"SMTP_USERNAME", "SMTP_USER"}, ""),
		SMTPPassword:                readEnvAny([]string{"SMTP_PASSWORD", "SMTP_PASS"}, ""),
		SMTPFromEmail:               readEnvAny([]string{"SMTP_FROM_EMAIL", "EMAIL_FROM"}, ""),
		SMTPFromName:                readEnvAny([]string{"SMTP_FROM_NAME", "EMAIL_FROM_NAME"}, "Mapa do Corre"),
		EmailCodeTTLMinutes:         readEnvInt("EMAIL_CODE_TTL_MINUTES", 10),
		EmailCodeMaxAttempts:        readEnvInt("EMAIL_CODE_MAX_ATTEMPTS", 5),
		WriteRateLimitMax:           readEnvInt("WRITE_RATE_LIMIT_MAX", 20),
		WriteRateLimitWindowSeconds: readEnvInt("WRITE_RATE_LIMIT_WINDOW_SECONDS", 60),
	}

	if appConfig.Port == "" {
		return Config{}, fmt.Errorf("PORT nao pode ser vazio")
	}

	if appConfig.FrontendOrigin == "" {
		return Config{}, fmt.Errorf("FRONTEND_ORIGIN nao pode ser vazio")
	}

	if appConfig.PersistenceMode != "memory" && appConfig.PersistenceMode != "postgres" {
		return Config{}, fmt.Errorf("PERSISTENCE_MODE precisa ser 'memory' ou 'postgres'")
	}

	if appConfig.EmailCodeTTLMinutes <= 0 {
		return Config{}, fmt.Errorf("EMAIL_CODE_TTL_MINUTES precisa ser maior que zero")
	}

	if appConfig.EmailCodeMaxAttempts <= 0 {
		return Config{}, fmt.Errorf("EMAIL_CODE_MAX_ATTEMPTS precisa ser maior que zero")
	}

	if appConfig.WriteRateLimitMax <= 0 {
		return Config{}, fmt.Errorf("WRITE_RATE_LIMIT_MAX precisa ser maior que zero")
	}

	if appConfig.WriteRateLimitWindowSeconds <= 0 {
		return Config{}, fmt.Errorf("WRITE_RATE_LIMIT_WINDOW_SECONDS precisa ser maior que zero")
	}

	return appConfig, nil
}

func (config Config) HasSMTPConfiguration() bool {
	return strings.TrimSpace(config.SMTPHost) != "" &&
		strings.TrimSpace(config.SMTPPort) != "" &&
		strings.TrimSpace(config.SMTPUsername) != "" &&
		strings.TrimSpace(config.SMTPPassword) != "" &&
		strings.TrimSpace(config.SMTPFromEmail) != ""
}

func loadEnvFiles() {
	// Mantemos a ordem do mais proximo para o mais abrangente.
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load("../.env.local")
}

func readEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		return value
	}

	return fallback
}

func readEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}

	return fallback
}

func readEnvInt(key string, fallback int) int {
	rawValue := strings.TrimSpace(os.Getenv(key))
	if rawValue == "" {
		return fallback
	}

	parsedValue, err := strconv.Atoi(rawValue)
	if err != nil {
		return fallback
	}

	return parsedValue
}
