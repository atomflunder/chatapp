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

func main() {
	var username string

	fmt.Println("Hi, please type in your Username: ")
	fmt.Scan(&username)

	cfg := models.GetConfig()

	go sendLoop(cfg, username)
	writeLoop(cfg, username)

}

func sendLoop(cfg models.Config, username string) {
	var secondsSleep time.Duration = 2

	for {
		timestamp := time.Now().UnixMilli() - int64(1000*secondsSleep)

		var newMsgs []models.Message

		resp, err := http.Get(fmt.Sprintf("http://%s:%s/messages?Since=%d", cfg.Host, cfg.Port, timestamp))
		if err != nil {
			fmt.Println("Could not get new messages")
			time.Sleep(time.Second * secondsSleep * 2)
		}

		err = json.NewDecoder(resp.Body).Decode(&newMsgs)
		if err != nil {
			fmt.Println(err, resp.Body)
			time.Sleep(time.Second * secondsSleep * 2)
		}

		for _, msg := range newMsgs {
			fmt.Printf("\n%s", msg.Format())
		}
		if len(newMsgs) > 0 {
			fmt.Printf("\n%s (%s): ", username, time.Now().Format(time.TimeOnly))
		}

		resp.Body.Close()

		time.Sleep(time.Second * secondsSleep)
	}
}

func writeLoop(cfg models.Config, username string) {
	for {
		var content string

		fmt.Printf("%s (%s): ", username, time.Now().Format(time.TimeOnly))
		fmt.Scanf("%s\n", &content)

		part := models.ParialMessage{Username: username, Content: content}
		message := part.GetMessage()

		messageJSON, err := json.Marshal(message)

		if err != nil {
			log.Fatal("Could not marshal message content")
		}

		resp, err := http.Post(fmt.Sprintf("http://%s:%s/messages/new", cfg.Host, cfg.Port), "application/json", bytes.NewBuffer([]byte(messageJSON)))
		if err != nil {
			log.Fatal("Error sending request!")
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
