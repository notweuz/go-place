package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func Logger(ctx fiber.Ctx) error {
	log.Debug().Str("ip", ctx.IP()).Msgf("New HTTP %s %s request", ctx.Method(), ctx.Path())
	start := time.Now()

	err := ctx.Next()
	latency := time.Since(start)

	if err != nil {
		return err
	}

	log.Info().Msgf("HTTP %s %s completed with status %d in %v", ctx.Method(), ctx.Path(), ctx.Response().StatusCode(), latency)
	return nil
}
