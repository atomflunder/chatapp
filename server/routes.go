package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/atomflunder/chatapp/models"
)

type Handler struct {
	w *DBWrapper
}

func newHandler(w *DBWrapper) *Handler {
	return &Handler{w}
}

func (h *Handler) registerRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /messages", h.getMessages)
	r.HandleFunc("POST /messages/new", h.postMessage)

	return r
}

func (h *Handler) getMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("New GET Messages Request received from %s\n", r.RemoteAddr)

	query := r.URL.Query()
	username := query.Get("Username")
	channel := query.Get("Channel")
	since := query.Get("Since")

	time, err := strconv.Atoi(since)
	if err != nil {
		time = 0
	}

	var ms []models.Message = h.w.getMessages(channel, username, int64(time))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms)
}

func (h *Handler) postMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("New POST Messages Request received from %s\n", r.RemoteAddr)

	var m models.ParialMessage

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

	if m.Content == "" {
		return
	}

	h.w.insertMessage(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(m)
}
