package main

import (
	"log"

	"github.com/atomflunder/chatapp/database"
)

func main() {
	w, err := database.OpenDB()
	if err != nil {
		log.Fatal("Error opening database")
	}
	defer w.Db.Close()

	w.Initialize()

	InitializeRoutes(w)
}
