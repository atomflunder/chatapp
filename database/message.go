package database

import (
	"fmt"
	"time"
)

type Message struct {
	ID        string
	Content   string
	Timestamp int64
	Username  string
}

type ParialMessage struct {
	Content  string
	Username string
}

func (m Message) Format() string {
	return fmt.Sprintf("%s (%s): %s\n", m.Username, time.Unix(0, m.Timestamp*int64(time.Millisecond)).Format(time.TimeOnly), m.Content)
}
