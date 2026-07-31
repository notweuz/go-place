package config

import (
	"go-place/internal/errs"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DatabaseDSN string
	LogLevel    string
	JwtSecret   string
	AppPort     int
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file, using environment variables")
	}

	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		log.Warn().Msg("Database DSN not set, fallback to default")
		databaseDSN = "host=localhost user=postgres password=postgres dbname=go-place port=5432 sslmode=disable"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		log.Warn().Msg("Log Level not set, fallback to default")
		logLevel = "info"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errs.ErrNoJWTSecretFound
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Warn().Msgf("Failed to parse APP_PORT")
		appPort = 8080
	}

	cfg.DatabaseDSN = databaseDSN
	cfg.LogLevel = logLevel
	cfg.JwtSecret = jwtSecret
	cfg.AppPort = appPort

	log.Info().Msg("Config loaded successfully!")
	return cfg, nil
}
