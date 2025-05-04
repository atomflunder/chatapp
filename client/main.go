package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/atomflunder/chatapp/models"
)

func main() {
	var username string

	// TODO: Maybe also remove these scans for bufio reader
	fmt.Println("Hi, please type in your Username: ")
	fmt.Scan(&username)

	var channel string

	fmt.Println("Type in the channel you want to connect to: ")
	fmt.Scan(&channel)

	fmt.Printf("Connecting to %s\n", channel)

	cfg := models.GetConfig()

	go sendLoop(cfg, username, channel)
	writeLoop(cfg, username, channel)

}

func sendLoop(cfg models.Config, username string, channel string) {
	var secondsSleep time.Duration = 2

	for {
		timestamp := time.Now().UnixMilli() - int64(1000*secondsSleep)

		var newMsgs []models.Message

		resp, err := http.Get(fmt.Sprintf("http://%s:%s/messages?Since=%d&Channel=%s&Username=%s", cfg.Host, cfg.Port, timestamp, channel, username))
		if err != nil {
			fmt.Println("Could not get new messages")
		}

		err = json.NewDecoder(resp.Body).Decode(&newMsgs)
		if err != nil {
			fmt.Println(err, resp.Body)
		}

		// TODO: Make your already typed in message not disappear when another user sends one
		if len(newMsgs) > 0 {
			fmt.Printf("\033[1A\033[K") // Deletes the last line
			for _, msg := range newMsgs {
				fmt.Printf("\n%s", msg.Format())
			}
			fmt.Printf("\n%s >: ", username)
		}

		resp.Body.Close()

		time.Sleep(time.Second * secondsSleep)
	}
}

func writeLoop(cfg models.Config, username string, channel string) {
	inputReader := bufio.NewReader(os.Stdin)

	for {
		var content string

		fmt.Printf("%s >: ", username)

		content, err := inputReader.ReadString('\n')
		if err != nil {
			log.Fatal("Could not read input!")
		}

		content = strings.TrimSuffix(content, "\n")

		part := models.PartialMessage{Username: username, Content: content}
		message := part.GetMessage(channel)

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
