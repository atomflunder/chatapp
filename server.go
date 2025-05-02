package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func StartServer(db *sql.DB) {
	handler := NewHandler(db)
	messageRouter := handler.RegisterRoutes()

	router := http.NewServeMux()
	router.Handle("/", messageRouter)

	fmt.Println("Server up and running!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
