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
			cl := h.getClientsInChannel(client.identity.Channel)
			h.sendSystemMessage(client.identity.Channel, fmt.Sprintf("User %s joined, there are now %d user(s) in this chat",
				client.identity.Username, len(cl)+1))

			if h.isIdentityInUse(client.identity) {
				partialMessage := models.PartialMessage{
					Content: "Username/Channel combination already in use, please reconnect!",
					Identity: models.Identity{
						Username: "system",
						Channel:  client.identity.Channel,
					},
				}
				errorMessage, err := json.Marshal(partialMessage.GetMessage())
				if err == nil {
					client.send <- errorMessage
				}
				close(client.send)
				continue
			}

			clNames := h.getAllClientNames(client.identity.Channel)

			h.clients[client] = true
			h.sendPrivateMessage(client, fmt.Sprintf("Welcome to #%s, there are %d other user(s) in this chat: %s",
				client.identity.Channel, len(cl), clNames))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			cl := h.getClientsInChannel(client.identity.Channel)
			h.sendSystemMessage(client.identity.Channel, fmt.Sprintf("User %s left, there are now %d user(s) in this chat",
				client.identity.Username, len(cl)))
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

// Sends a system message in a channel.
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

// Sends a private message directly to the client.
func (h *Hub) sendPrivateMessage(client *Client, content string) {
	partial := models.PartialMessage{
		Content: content,
		Identity: models.Identity{
			Username: "system",
			Channel:  "",
		},
	}
	sysMsg := partial.GetMessage()
	byteMsg, err := json.Marshal(sysMsg)
	if err != nil {
		return
	}

	client.send <- byteMsg
}

// Gets you every client in a channel.
func (h *Hub) getClientsInChannel(channel string) []models.Identity {
	clients := []models.Identity{}

	for client := range h.clients {
		if client.identity.Channel == channel {
			clients = append(clients, client.identity)
		}
	}

	return clients
}

// Gets every client's username, separated by spaces.
func (h *Hub) getAllClientNames(channel string) string {
	s := ""
	cl := h.getClientsInChannel(channel)

	for _, client := range cl {
		s += fmt.Sprintf("%s ", client.Username)
	}

	return s
}

// Checks if a username/channel combination is already connected to the hub.
func (h *Hub) isIdentityInUse(identity models.Identity) bool {
	cl := h.getClientsInChannel(identity.Channel)

	for _, client := range cl {
		if client.Channel == identity.Channel && client.Username == identity.Username {
			return true
		}
	}

	return false
}
