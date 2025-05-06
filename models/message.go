package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Username  string `json:"username"`
	Channel   string `json:"channel"`
}

type PartialMessage struct {
	Content  string `json:"content"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

// Builds a full message from a partial one.
func (p PartialMessage) GetMessage() Message {
	return Message{
		ID:        uuid.NewString(),
		Timestamp: time.Now().UnixMilli(),
		Username:  p.Username,
		Content:   p.Content,
		Channel:   p.Channel,
	}
}

// Formats a message in a nice readable format.
func (m Message) Format() string {
	return fmt.Sprintf("%s (%s):\n%s", m.Username, time.Unix(0, m.Timestamp*int64(time.Millisecond)).Format(time.TimeOnly), m.Content)
}
