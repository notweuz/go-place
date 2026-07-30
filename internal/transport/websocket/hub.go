package websocket

import (
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
}

func NewHub() Hub {
	return &hub{
		clients:    make(map[Client]bool),
		register:   make(chan Client),
		unregister: make(chan Client),
		broadcast:  make(chan message.Base),
	}
}

func (h *hub) Broadcast(msg message.Base) {
	log.Info().Interface("message", msg).Msg("Broadcasting message to all clients")
	h.broadcast <- msg
	log.Info().Msg("Broadcasting complete")
}

func (h *hub) Register(client Client) {
	log.Info().Msg("Registering client")
	h.register <- client
}

func (h *hub) Unregister(client Client) {
	log.Info().Msg("Unregistering client")
	h.unregister <- client
}

func (h *hub) Run() {
	for {
		select {
		case cl := <-h.register:
			if _, ok := h.clients[cl]; !ok {
				log.Info().Msg("New websocket client registered")
				h.clients[cl] = true
			}
		case cl := <-h.unregister:
			if _, ok := h.clients[cl]; ok {
				cl.Close()
				delete(h.clients, cl)
			}
		case msg := <-h.broadcast:
			for cl := range h.clients {
				cl.Send(msg)
			}
		}
	}
}
