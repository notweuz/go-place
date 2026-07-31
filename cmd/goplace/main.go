package main

import (
	"go-place/internal/config"
	"go-place/internal/logger"
	"go-place/internal/server"
	"os"
	"os/signal"

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
	app.Start()

	log.Debug().Msg("Initialized graceful shutdown")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	app.Shutdown()
}
