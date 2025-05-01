package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Message struct {
	ID        int
	Content   string
	Timestamp int
	Username  string
}

var messages = make(map[int]Message)
var currentID = 1

func messageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getMessages(w, r)
	case "POST":
		postMessage(w, r)
	default:
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	ms := make([]Message, 0, len(messages))

	query := r.URL.Query()
	username := query.Get("Username")

	for _, m := range messages {
		if username == "" || m.Username == username {
			ms = append(ms, m)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms)
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var m Message

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading Body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	m.ID = currentID
	m.Timestamp = 123
	currentID++

	messages[m.ID] = m
	fmt.Println("New Message Received - ID", m.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(m)
}

func StartServer() {
	http.HandleFunc("/new", messageHandler)

	fmt.Println("Server up and running!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
