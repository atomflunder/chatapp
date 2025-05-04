package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

	p := tea.NewProgram(initialModel(username, channel))

	go fetchNewMessages(p)

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

// Sends the updateMessage to the tea.Update function, which triggers a re-fetch of the messages since the last update.
func fetchNewMessages(p *tea.Program) {
	var secondsSleep time.Duration = 2

	for {
		p.Send(updateMessage{
			lastUpdate: time.Now().UnixMilli() - (int64(secondsSleep * 1000)),
		})
		time.Sleep(time.Second * secondsSleep)
	}
}
