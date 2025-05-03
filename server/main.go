package main

import (
	"log"
)

func main() {
	w, err := OpenDB()
	if err != nil {
		log.Fatal("Error opening database")
	}
	defer w.Db.Close()

	w.initialize()

	InitializeRoutes(w)
}
