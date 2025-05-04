package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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

	resp, err := http.Get(fmt.Sprintf("http://%s:%s/channels/%s/messages?Since=%d&Username=%s", cfg.Host, cfg.Port, channel, timestamp, username))
	if err != nil {
		fmt.Println("Could not get new messages")
		return []models.Message{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Received status %d\n", resp.StatusCode)
		return []models.Message{}
	}

	var newMsgs []models.Message
	err = json.NewDecoder(resp.Body).Decode(&newMsgs)
	if err != nil {
		fmt.Println(err, resp.Body)
		return []models.Message{}
	}

	return newMsgs
}
