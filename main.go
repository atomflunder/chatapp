package main

import "log"

func main() {
	w, err := openDB()
	if err != nil {
		log.Fatal("Error opening database")
	}
	defer w.db.Close()

	w.Initialize()

	StartServer(w)
}
