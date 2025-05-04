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
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Hi, please type in your Username: ")
	username, err := getInput(inputReader)
	if err != nil {
		log.Fatal("Could not read input!")
	}

	fmt.Println("Type in the channel you want to connect to: ")
	channel, err := getInput(inputReader)
	if err != nil {
		log.Fatal("Could not read input!")
	}

	if strings.Contains(channel, " ") {
		log.Fatal("Channel IDs cannot have spaces in them")
	}

	fmt.Printf("Connecting to %s\n", channel)

	cfg := models.GetConfig()

	go sendLoop(cfg, username, channel)
	writeLoop(cfg, username, channel, inputReader)

}

func getInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	input = strings.TrimSuffix(input, "\n")

	return input, nil
}

func sendLoop(cfg models.Config, username string, channel string) {
	var secondsSleep time.Duration = 2

	for {
		timestamp := time.Now().UnixMilli() - int64(1000*secondsSleep)

		var newMsgs []models.Message

		resp, err := http.Get(fmt.Sprintf("http://%s:%s/channels/%s/messages?Since=%d&Username=%s", cfg.Host, cfg.Port, channel, timestamp, username))
		if err != nil {
			fmt.Println("Could not get new messages")
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Received status %d\n", resp.StatusCode)
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

func writeLoop(cfg models.Config, username string, channel string, reader *bufio.Reader) {
	for {
		fmt.Printf("%s >: ", username)

		content, err := getInput(reader)
		if err != nil {
			log.Fatal("Could not read input!")
		}

		part := models.PartialMessage{Username: username, Content: content}
		message := part.GetMessage(channel)

		messageJSON, err := json.Marshal(message)

		if err != nil {
			log.Fatal("Could not marshal message content")
		}

		resp, err := http.Post(fmt.Sprintf("http://%s:%s/channels/%s/messages", cfg.Host, cfg.Port, channel), "application/json", bytes.NewBuffer([]byte(messageJSON)))
		if err != nil {
			log.Fatal("Error sending request!")
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
