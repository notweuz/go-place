package websocket

import (
	"go-place/internal/transport/websocket/message"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	Upgrade() fiber.Handler
}

type handler struct {
	hub Hub
}

func NewHandler(hub Hub) Handler {
	return &handler{hub: hub}
}

func (h *handler) Upgrade() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		cl := NewClient(c)

		h.hub.Register(cl)
		defer h.hub.Unregister(cl)

		go cl.Write()

		for {
			var msg message.Base
			err := c.ReadJSON(&msg)
			if err != nil {
				break
			}
		}
	})
}
