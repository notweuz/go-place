package websocket

import (
	"context"
	"go-place/internal/transport/websocket/message"

	"github.com/rs/zerolog/log"
)

type Hub interface {
	Run()
	Broadcast(message message.Base)
	Register(Client)
	Unregister(Client)
}

type hub struct {
	clients    map[Client]bool
	register   chan Client
	unregister chan Client
	broadcast  chan message.Base
	ctx        context.Context
}

func NewHub(ctx context.Context) Hub {
	return &hub{
		clients:    make(map[Client]bool),
		register:   make(chan Client),
		unregister: make(chan Client),
		broadcast:  make(chan message.Base),
		ctx:        ctx,
	}
}

func (h *hub) Broadcast(msg message.Base) {
	select {
	case h.broadcast <- msg:
		log.Info().Interface("message", msg).Msg("Broadcasting message")
	case <-h.ctx.Done():
		log.Warn().Msg("Hub is shut down, broadcast skipped")
	}
}
func (h *hub) Register(client Client) {
	select {
	case h.register <- client:
		log.Info().Msg("Registering client")
	case <-h.ctx.Done():
		log.Warn().Msg("Hub is shut down, register skipped")
	}
}
func (h *hub) Unregister(client Client) {
	select {
	case h.unregister <- client:
		log.Info().Msg("Unregistering client")
	case <-h.ctx.Done():
		log.Warn().Msg("Hub is shut down, unregister skipped")
	}
}

func (h *hub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			for cl := range h.clients {
				cl.Close()
				delete(h.clients, cl)
			}
			return
		case cl := <-h.register:
			if _, ok := h.clients[cl]; !ok {
				log.Info().Msg("New websocket client registered")
				h.clients[cl] = true
			}
		case cl := <-h.unregister:
			if _, ok := h.clients[cl]; ok {
				delete(h.clients, cl)
			}
		case msg := <-h.broadcast:
			for cl := range h.clients {
				cl.Send(msg)
			}
		}
	}
}
