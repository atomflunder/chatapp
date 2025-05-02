package main

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer(w *DBWrapper) {
	handler := NewHandler(w)
	messageRouter := handler.RegisterRoutes()

	router := http.NewServeMux()
	router.Handle("/", messageRouter)

	fmt.Println("Server up and running!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
