package middleware

import (
	"go-place/internal/auth"
	"go-place/internal/config"
	"go-place/internal/errs"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func AuthProtected(cfg *config.Config) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
			return errs.UnauthorizedTokenError(errs.ErrNotAuthorized, "token is missing")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := auth.ParseToken(tokenString, cfg.JwtSecret)
		if err != nil {
			return errs.UnauthorizedTokenError(err, "token is invalid")
		}

		ctx.Locals("user_id", userID)

		return ctx.Next()
	}
}

func GetCurrentUserID(ctx fiber.Ctx) (uint64, error) {
	userID, ok := ctx.Locals("user_id").(uint64)
	if !ok {
		log.Error().Msg("failed to get user id in context")
		return 0, errs.UnauthorizedTokenError(errs.ErrNotAuthorized, "failed to get user id in context")
	}
	return userID, nil
}
