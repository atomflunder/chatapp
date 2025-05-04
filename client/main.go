package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	p := tea.NewProgram(initialModel(username, channel))

	go sendLoop(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func getInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	input = strings.TrimSuffix(input, "\n")

	return input, nil
}
