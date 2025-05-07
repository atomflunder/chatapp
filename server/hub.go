package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/atomflunder/chatapp/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	dbWrapper  *DBWrapper
}

func newHub() *Hub {
	w, err := OpenDB()
	if err != nil {
		log.Fatal("Error opening database")
	}

	w.initialize()

	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		dbWrapper:  w,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case rawMessage := <-h.broadcast:
			h.handleMessage(rawMessage)
		}
	}
}

// Parses, stores, and sends an incoming message to every client in the same channel, except the author.
func (h *Hub) handleMessage(rawMessage []byte) {
	partialMessage := models.PartialMessage{}
	err := json.Unmarshal(rawMessage, &partialMessage)
	if err != nil {
		log.Println("Error parsing message data")
		return
	}

	fmt.Printf("Received message from %s in %s\n", partialMessage.Username, partialMessage.Channel)

	h.dbWrapper.insertMessage(partialMessage)
	fullMessage := partialMessage.GetMessage()

	encodedMessage, err := json.Marshal(fullMessage)
	if err != nil {
		log.Println("Error marshalling message data")
		return
	}

	for client := range h.clients {
		if client.identity.Channel != fullMessage.Identity.Channel || client.identity.Username == fullMessage.Username {
			continue
		}

		select {
		case client.send <- encodedMessage:

		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
