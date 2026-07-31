package main

import (
	"fmt"
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
	go func() {
		log.Info().Msg("Starting fiber application")
		if err := app.Fiber.Listen(fmt.Sprintf(":%d", app.Config.AppPort)); err != nil {
			log.Panic().Err(err).Msg("Error starting app")
		}
	}()

	log.Debug().Msg("Initialized graceful shutdown")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	app.Shutdown()
}
