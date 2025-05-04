package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/atomflunder/chatapp/models"
)

func postMessage(message models.Message) {
	cfg := models.GetConfig()

	messageJSON, err := json.Marshal(message)

	if err != nil {
		log.Fatal("Could not marshal message contents")
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:%s/channels/%s/messages", cfg.Host, cfg.Port, message.Channel), "application/json", bytes.NewBuffer([]byte(messageJSON)))
	if err != nil {
		log.Fatal("Error sending request!")
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

}

func getMessages(username string, channel string, timestamp int64) []models.Message {
	cfg := models.GetConfig()

	var newMsgs []models.Message

	resp, err := http.Get(fmt.Sprintf("http://%s:%s/channels/%s/messages?Since=%d&Username=%s", cfg.Host, cfg.Port, channel, timestamp, username))
	if err != nil {
		newMsgs = []models.Message{{
			ID:        "0",
			Content:   "Could not receive new messages",
			Timestamp: time.Now().UnixMilli(),
			Username:  "system",
			Channel:   channel,
		},
		}
		return newMsgs
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		newMsgs = []models.Message{{
			ID:        "0",
			Content:   fmt.Sprintf("Received %d Status Code from system", resp.StatusCode),
			Timestamp: time.Now().UnixMilli(),
			Username:  "system",
			Channel:   channel,
		},
		}
		return newMsgs
	}

	err = json.NewDecoder(resp.Body).Decode(&newMsgs)
	if err != nil {
		newMsgs = []models.Message{{
			ID:        "0",
			Content:   "Could not decode new messages from server",
			Timestamp: time.Now().UnixMilli(),
			Username:  "system",
			Channel:   channel,
		},
		}
		return newMsgs
	}

	return newMsgs
}
