package main

import (
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	ID        string
	Content   string
	Timestamp int64
	Username  string
}

func StartServer() {
	handler := NewHandler()
	messageRouter := handler.RegisterRoutes()

	router := http.NewServeMux()
	router.Handle("/", messageRouter)

	fmt.Println("Server up and running!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
