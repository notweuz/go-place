package websocket

import (
	"go-place/internal/transport/websocket/message"
	"sync"

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
	done chan struct{}
	once sync.Once
}

func NewClient(conn *websocket.Conn) Client {
	return &client{
		conn: conn,
		send: make(chan message.Base, 256),
		done: make(chan struct{}),
	}
}

func (c *client) Send(msg message.Base) {
	select {
	case <-c.done:
		return
	case c.send <- msg:
		log.Debug().Interface("message", msg).Msg("Sending message to client")
	default:
		log.Warn().Msg("Client is not ready to receive message, aborting connection")
		c.Close()
	}
}

func (c *client) Write() {
	for {
		select {
		case msg := <-c.send:
			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Error().Err(err).Msg("Error writing to client")
				c.Close()
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *client) Close() {
	c.once.Do(func() {
		log.Debug().Msg("Closing connection with client")
		err := c.conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close connection with client")
		}
		close(c.done)
		close(c.send)
	})
}
