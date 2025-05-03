package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/atomflunder/chatapp/models"
	_ "github.com/mattn/go-sqlite3"
)

type DBWrapper struct{ Db *sql.DB }

func OpenDB() (*DBWrapper, error) {
	db, err := sql.Open("sqlite3", "./database/messages.db")

	if err != nil {
		log.Fatal("Error opening db")
		return nil, err
	}

	return &DBWrapper{db}, nil
}

func (w *DBWrapper) Initialize() {
	sqlStmt := `
	create table if not exists messages (id text not null primary key, username text, timestamp integer, content text);
	`
	_, err := w.Db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (w *DBWrapper) GetMessages() []models.Message {
	messages := []models.Message{}

	rows, err := w.Db.Query(`select * from messages`)
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

		message := models.Message{ID: id, Username: username, Timestamp: timestamp, Content: content}

		messages = append(messages, message)
	}

	return messages
}

func (w *DBWrapper) GetMessagesFromUser(u string) []models.Message {
	messages := []models.Message{}

	tx, err := w.Db.Begin()
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

		message := models.Message{ID: id, Username: username, Timestamp: timestamp, Content: content}

		messages = append(messages, message)
	}

	return messages
}

func (w *DBWrapper) InsertMessage(m models.ParialMessage) {
	tx, err := w.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into messages(id, username, timestamp, content) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	message := m.GetMessage()

	_, err = stmt.Exec(message.ID, message.Username, message.Timestamp, message.Content)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got new message", message.ID)
}
