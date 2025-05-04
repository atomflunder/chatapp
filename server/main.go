package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atomflunder/chatapp/models"
)

func main() {
	w, err := OpenDB()
	if err != nil {
		log.Fatal("Error opening database")
	}
	defer w.Db.Close()

	w.initialize()

	handler := newHandler(w)
	messageRouter := handler.registerRoutes()

	router := http.NewServeMux()
	router.Handle("/", messageRouter)

	config := models.GetConfig()

	fmt.Println("Server up and running!")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host, config.Port), router))
}
