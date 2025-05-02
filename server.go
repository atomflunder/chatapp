package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	Content   string
	Timestamp int64
	Username  string
}

var messages = make(map[uuid.UUID]Message)

func StartServer() {
	handler := NewHandler()
	messageRouter := handler.RegisterRoutes()

	router := http.NewServeMux()
	router.Handle("/", messageRouter)

	fmt.Println("Server up and running!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
