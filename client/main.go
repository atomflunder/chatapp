package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atomflunder/chatapp/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

func main() {
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Hi, please type in your Username: ")
	username, err := getInput(inputReader)
	if err != nil {
		log.Fatal("Could not read input!")
	}

	if username == "system" || strings.Contains(username, " ") {
		log.Fatal("Invalid username")
	}

	fmt.Println("Type in the channel you want to connect to: ")
	channel, err := getInput(inputReader)
	if err != nil {
		log.Fatal("Could not read input!")
	}

	if strings.Contains(channel, " ") {
		log.Fatal("Channel IDs cannot have spaces in them")
	}

	cfg := models.GetConfig()
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/channels/%s/user/%s", cfg.Host, cfg.Port, channel, username), nil)

	if err != nil {
		log.Fatal("Failed to connect to websocket")
	}
	defer c.Close()

	p := tea.NewProgram(initialModel(username, channel, c))

	go fetchNewMessages(p, c)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

// Gets the user input in a line and cuts off the new line character at the end.
func getInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	input = strings.TrimSuffix(input, "\n")

	return input, nil
}

func fetchNewMessages(p *tea.Program, c *websocket.Conn) {
	done := make(chan struct{})
	defer close(done)
	for {
		message := models.Message{}
		err := c.ReadJSON(&message)
		if err != nil {
			return
		}
		p.Send(newMessage{
			message: message,
		})
	}
}
