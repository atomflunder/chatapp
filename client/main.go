package main

import (
	"fmt"
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

	var username string

	fmt.Println("Hi, please type in your Username: ")
	fmt.Scan(&username)

	for {
		var content string

		fmt.Printf("%s: ", username)
		fmt.Scanf("%s\n", &content)

		w.InsertMessage(database.ParialMessage{Username: username, Content: content})
	}
}
