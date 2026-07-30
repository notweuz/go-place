package message

import "encoding/json"

type Base struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}
