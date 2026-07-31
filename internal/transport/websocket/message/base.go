package message

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type Base struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

func NewMessage(event string, payload interface{}) *Base {
	rawPayload, err := json.Marshal(payload)

	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal payload")
		return &Base{
			Event:   event,
			Payload: nil,
		}
	}

	return &Base{
		Event:   event,
		Payload: rawPayload,
	}
}
