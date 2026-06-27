package middleware

import (
	"go-place/internal/auth"
	"go-place/internal/config"
	"go-place/internal/errs"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthProtected(cfg *config.Config) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
			return errs.NewAppError(errs.ErrNotAuthorized, fiber.StatusUnauthorized, "auth token is missing")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := auth.ParseToken(tokenString, cfg.JwtSecret)
		if err != nil {
			return errs.NewAppError(errs.ErrNotAuthorized, fiber.StatusUnauthorized, "invalid token")
		}

		ctx.Locals("user_id", userID)

		return ctx.Next()
	}
}
