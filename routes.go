package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /messages", getMessages)
	r.HandleFunc("POST /messages/new", postMessage)

	return r
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

	uuid, err := uuid.NewUUID()
	if err != nil {
		http.Error(w, "Error Setting UUID", http.StatusInternalServerError)
		return
	}

	m.ID = uuid
	m.Timestamp = time.Now().Unix()

	messages[m.ID] = m
	fmt.Println("New Message Received - ID", m.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(m)
}
