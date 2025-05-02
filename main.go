package main

import "log"

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal("Error opening database")
	}
	defer db.Close()

	InitializeDB(db)

	StartServer(db)
}
