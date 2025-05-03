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
}

type ParialMessage struct {
	Content  string
	Username string
}

func (p ParialMessage) GetMessage() Message {
	return Message{
		ID:        uuid.NewString(),
		Timestamp: time.Now().UnixMilli(),
		Username:  p.Username,
		Content:   p.Content,
	}
}

func (m Message) Format() string {
	return fmt.Sprintf("%s (%s): %s\n", m.Username, time.Unix(0, m.Timestamp*int64(time.Millisecond)).Format(time.TimeOnly), m.Content)
}
