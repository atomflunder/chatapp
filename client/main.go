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

	for {
		var content string

		fmt.Printf("%s (%s): ", username, time.Now().Format(time.TimeOnly))
		fmt.Scanf("%s\n", &content)

		// TODO: Implement getting messages
		// 		 Implement multiple connections, and distinguishing between them

		part := models.ParialMessage{Username: username, Content: content}
		message := part.GetMessage()

		messageJSON, err := json.Marshal(message)

		if err != nil {
			log.Fatal("Could not marshal message content")
		}

		config := models.GetConfig()

		resp, err := http.Post(fmt.Sprintf("http://%s:%s/messages/new", config.Host, config.Port), "application/json", bytes.NewBuffer([]byte(messageJSON)))
		if err != nil {
			log.Fatal("Error sending request!")
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
