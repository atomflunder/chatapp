package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atomflunder/chatapp/database"
)

func InitializeRoutes(w *database.DBWrapper) {
	handler := NewHandler(w)
	messageRouter := handler.RegisterRoutes()

	router := http.NewServeMux()
	router.Handle("/", messageRouter)

	fmt.Println("Server up and running!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
