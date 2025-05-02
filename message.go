package main

import "fmt"

type Message struct {
	ID        string
	Content   string
	Timestamp int64
	Username  string
}

func (m Message) FormatMessage() string {
	return fmt.Sprintf("Message %s by %s\nContent: %s\n", m.ID, m.Username, m.Content)
}
