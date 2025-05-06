package models

import (
	"strings"
	"testing"
	"time"
)

func TestPartialMessage_GetMessage(t *testing.T) {
	p := PartialMessage{
		Content:  "Hello, world!",
		Username: "alice",
		Channel:  "general",
	}
	channel := "general"
	msg := p.GetMessage()

	if msg.Content != p.Content {
		t.Errorf("Expected content %q, got %q", p.Content, msg.Content)
	}
	if msg.Username != p.Username {
		t.Errorf("Expected username %q, got %q", p.Username, msg.Username)
	}
	if msg.Channel != channel {
		t.Errorf("Expected channel %q, got %q", channel, msg.Channel)
	}
	if msg.ID == "" {
		t.Error("Expected a non-empty ID")
	}
	if msg.Timestamp == 0 {
		t.Error("Expected a non-zero timestamp")
	}
}

func TestMessage_Format(t *testing.T) {
	now := time.Now()
	msg := Message{
		ID:        "test-id",
		Content:   "Test message",
		Username:  "bob",
		Channel:   "random",
		Timestamp: now.UnixMilli(),
	}

	formatted := msg.Format()

	if !strings.Contains(formatted, "bob") {
		t.Errorf("Expected formatted string to contain username, got: %s", formatted)
	}
	if !strings.Contains(formatted, "Test message") {
		t.Errorf("Expected formatted string to contain content, got: %s", formatted)
	}

	formattedTime := time.UnixMilli(msg.Timestamp).Format(time.TimeOnly)
	if !strings.Contains(formatted, formattedTime) {
		t.Errorf("Expected formatted string to contain time %s, got: %s", formattedTime, formatted)
	}
}
