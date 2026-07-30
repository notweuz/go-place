package websocket

import (
	"go-place/internal/transport/websocket/message"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Send(msg message.Base)
	Write()
	Close()
}

type client struct {
	conn *websocket.Conn
	send chan message.Base
}

func NewClient(conn *websocket.Conn) Client {
	return &client{
		conn: conn,
		send: make(chan message.Base, 256),
	}
}

func (c *client) Send(msg message.Base) {
	log.Debug().Interface("message", msg).Msg("Sending message to client")
	c.send <- msg
}

func (c *client) Write() {
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Error().Err(err).Msg("Error writing to client")
			c.Close()
			return
		}
	}
}

func (c *client) Close() {
	log.Debug().Msg("Closing connection with client")
	err := c.conn.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close connection with client")
	}
}
