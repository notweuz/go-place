package main

import (
	"go-place/internal/config"
	"go-place/internal/logger"
	"go-place/internal/server"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.SetupLogger()
	log.Info().Msg("Starting server")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Error loading config")
	}
	logger.UpdateLogLevel(cfg.LogLevel)

	app := server.NewApp(cfg)
	_ = app
}
