package middleware

import (
	"go-place/internal/errs"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func WebsocketUpgrade(ctx fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		return ctx.Next()
	}
	return errs.UpdateRequired(errs.ErrUpgradeRequired)
}
