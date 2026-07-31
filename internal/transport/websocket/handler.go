package websocket

import (
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
		cl := NewClient(c, func(cli Client) {
			h.hub.Unregister(cli)
		})

		h.hub.Register(cl)
		defer cl.Close()

		go cl.Write()

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	})
}
