package network

import (
	"encoding/json"
)

type MessageType string

const (
	MessageHandshake   MessageType = "HANDSHAKE"
	MessageBlock       MessageType = "BLOCK"
	MessageTransaction MessageType = "TRANSACTION"
	MessageSync        MessageType = "SYNC"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload []byte      `json:"payload"`
}

func CreateMessage(
	messageType MessageType,
	data interface{},
) (*Message, error) {

	payload, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    messageType,
		Payload: payload,
	}, nil
}

func DecodeMessage(
	message *Message,
	output interface{},
) error {

	return json.Unmarshal(
		message.Payload,
		output,
	)
}
