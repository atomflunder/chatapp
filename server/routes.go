package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
	r.HandleFunc("/channels/", h.routeHandler)

	return r
}

// Handles routing.
// Needs to handle every route, because the channel URL param is dynamic.
func (h *Handler) routeHandler(w http.ResponseWriter, r *http.Request) {

	params := strings.Split(strings.TrimPrefix(r.URL.Path, "/channels/"), "/")

	if len(params) != 2 {
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	channel, m := params[0], params[1]

	if m != "messages" {
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		h.getMessages(w, r, channel)
		return
	case "POST":
		h.postMessage(w, r, channel)
		return
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
}

// Fetches the requested messages out of the database.
func (h *Handler) getMessages(w http.ResponseWriter, r *http.Request, channel string) {
	fmt.Printf("New GET Messages Request received from %s\n", r.RemoteAddr)

	query := r.URL.Query()
	username := query.Get("Username")
	since := query.Get("Since")

	time, err := strconv.Atoi(since)
	if err != nil {
		time = 0
	}

	var ms []models.Message = h.w.getMessages(channel, username, int64(time))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms)
}

// Writes a new message to the database.
func (h *Handler) postMessage(w http.ResponseWriter, r *http.Request, channel string) {
	fmt.Printf("New POST Messages Request received from %s\n", r.RemoteAddr)

	var m models.PartialMessage

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

	h.w.insertMessage(m, channel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(m)
}
