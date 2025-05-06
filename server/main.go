package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/atomflunder/chatapp/models"
)

func main() {
	hub := newHub()
	go hub.run()
	defer hub.dbWrapper.Db.Close()

	http.HandleFunc("/channels/", func(w http.ResponseWriter, r *http.Request) {
		params := strings.Split(strings.TrimPrefix(r.URL.Path, "/channels/"), "/")

		if len(params) != 3 {
			http.Error(w, "Route not found", http.StatusNotFound)
			return
		}

		channel, u, username := params[0], params[1], params[2]

		if u != "user" {
			return
		}

		serveWs(hub, w, r, username, channel)
	})

	cfg := models.GetConfig()

	fmt.Println("Server up and running")
	http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), nil)
}
