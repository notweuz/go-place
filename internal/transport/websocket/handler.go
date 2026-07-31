package websocket

import (
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

const (
	pongWait = 60 * time.Second
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

		go cl.Write()

		h.hub.Register(cl)
		defer h.hub.Unregister(cl)
		defer cl.Close()

		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// websocket is read-only right now, so ReadMessage is used only for connection loss detection purposes
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	})
}
