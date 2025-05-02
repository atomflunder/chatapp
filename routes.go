package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	w *DBWrapper
}

func NewHandler(w *DBWrapper) *Handler {
	return &Handler{w}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /messages", h.getMessages)
	r.HandleFunc("POST /messages/new", h.postMessage)

	return r
}

func (h *Handler) getMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	username := query.Get("Username")

	var ms []Message
	if username != "" {
		ms = h.w.GetMessagesFromUser(username)
	} else {
		ms = h.w.GetMessages()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms)
}

func (h *Handler) postMessage(w http.ResponseWriter, r *http.Request) {
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

	m.ID = uuid.String()
	m.Timestamp = time.Now().UnixMilli()

	fmt.Println(m.FormatMessage())
	h.w.InsertMessage(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(m)
}
