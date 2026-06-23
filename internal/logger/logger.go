package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func SetupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.LevelColors[zerolog.DebugLevel] = 35
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func UpdateLogLevel(new string) {
	level, err := zerolog.ParseLevel(new)
	if err != nil {
		log.Error().Str("lvl", new).Msg("Failed to parse log level")
		log.Warn().Msg("Fallback to default one")
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Info().Msgf("Setting log level to %s", level.String())
}
