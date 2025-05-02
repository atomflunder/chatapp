package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database/messages.db")

	if err != nil {
		log.Fatal("Error opening db")
		return nil, err
	}

	return db, nil
}

func CreateDB() {
	db, err := openDB()
	if err != nil {
		log.Fatal("Error opening db")
		return
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists messages (id text not null primary key, username text, timestamp integer, content text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func GetMessages() []Message {
	messages := []Message{}

	db, err := openDB()
	if err != nil {
		log.Fatal("Error opening db")
		return messages
	}
	defer db.Close()

	rows, err := db.Query(`select * from messages`)
	if err != nil {
		log.Fatal(err)
		return messages
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var username string
		var timestamp int64
		var content string

		err = rows.Scan(&id, &username, &timestamp, &content)
		if err != nil {
			log.Fatal(err)
			return messages
		}

		message := Message{ID: id, Username: username, Timestamp: timestamp, Content: content}

		messages = append(messages, message)
	}

	return messages
}

func GetMessagesFromUser(u string) []Message {
	messages := []Message{}

	db, err := openDB()
	if err != nil {
		log.Fatal("Error opening db")
		return messages
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return messages
	}

	stmt, err := tx.Prepare(`select * from messages where username = ?`)
	if err != nil {
		log.Fatal(err)
		return messages
	}

	rows, err := stmt.Query(u)
	if err != nil {
		log.Fatal(err)
		return messages
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var username string
		var timestamp int64
		var content string

		err = rows.Scan(&id, &username, &timestamp, &content)
		if err != nil {
			log.Fatal(err)
			return messages
		}

		message := Message{ID: id, Username: username, Timestamp: timestamp, Content: content}

		messages = append(messages, message)
	}

	db.Close()

	return messages
}

func InsertMessage(m Message) {
	db, err := openDB()
	if err != nil {
		log.Fatal("Error opening db")
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into messages(id, username, timestamp, content) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.ID, m.Username, m.Timestamp, m.Content)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Database is locked after 2 inserts
}
