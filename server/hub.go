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
			if h.isIdentityInUse(client.identity) {
				// TODO: Tell the client it couldn't connect
				close(client.send)
				continue
			}

			cl := h.getClientsInChannel(client.identity.Channel)
			h.sendSystemMessage(client.identity.Channel, fmt.Sprintf("User %s joined, there are now %d users in this chat", client.identity.Username, len(cl)+1))
			h.clients[client] = true
			h.sendPrivateMessage(client.identity, fmt.Sprintf("Welcome to %s, there are %d other users in this chat", client.identity.Channel, len(cl)))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			cl := h.getClientsInChannel(client.identity.Channel)
			h.sendSystemMessage(client.identity.Channel, fmt.Sprintf("User %s left, there are now %d users in this chat", client.identity.Username, len(cl)))
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

func (h *Hub) sendSystemMessage(channel string, content string) {
	partial := models.PartialMessage{
		Content: content,
		Identity: models.Identity{
			Username: "system",
			Channel:  channel,
		},
	}
	sysMsg := partial.GetMessage()
	byteMsg, err := json.Marshal(sysMsg)
	if err != nil {
		return
	}

	for client := range h.clients {
		if client.identity.Channel == channel {
			client.send <- byteMsg
		}
	}
}

func (h *Hub) sendPrivateMessage(identity models.Identity, content string) {
	partial := models.PartialMessage{
		Content: content,
		Identity: models.Identity{
			Username: "system",
			Channel:  identity.Channel,
		},
	}
	sysMsg := partial.GetMessage()
	byteMsg, err := json.Marshal(sysMsg)
	if err != nil {
		return
	}

	for client := range h.clients {
		if client.identity.Username == identity.Username && client.identity.Channel == identity.Channel {
			client.send <- byteMsg
		}
	}
}

func (h *Hub) getClientsInChannel(channel string) []models.Identity {
	clients := []models.Identity{}

	for client := range h.clients {
		if client.identity.Channel == channel {
			clients = append(clients, client.identity)
		}
	}

	return clients
}

func (h *Hub) isIdentityInUse(identity models.Identity) bool {
	cl := h.getClientsInChannel(identity.Channel)

	for _, client := range cl {
		if client.Channel == identity.Channel && client.Username == identity.Username {
			return true
		}
	}

	return false
}
