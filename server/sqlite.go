package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/atomflunder/chatapp/models"
	_ "github.com/mattn/go-sqlite3"
)

// Thin wrapper around a sql.DB struct, in order to add some methods.
type DBWrapper struct{ Db *sql.DB }

// Opens the Database and returns a thin wrapper around it.
// Remember to use `defer DBWrapper.Db.Close()`.
func OpenDB() (*DBWrapper, error) {
	db, err := sql.Open("sqlite3", "./database/messages.db")

	if err != nil {
		log.Fatal("Error opening db")
		return nil, err
	}

	return &DBWrapper{db}, nil
}

// Initialize the database with the needed structure, for first time setup.
func (w *DBWrapper) initialize() {
	sqlStmt := `
	create table if not exists messages (id text not null primary key, username text, timestamp integer, channel text, content text);
	`
	_, err := w.Db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// Gets you all messages in a channel, excluding a username, and after a timestamp.
// To get all messages in a channel, pass in "" and 0.
func (w *DBWrapper) getMessages(c string, u string, t int64) []models.Message {
	messages := []models.Message{}

	tx, err := w.Db.Begin()
	if err != nil {
		log.Fatal(err)
		return messages
	}

	var rows *sql.Rows

	stmt, err := tx.Prepare(`select * from messages where timestamp >= ? and channel = ? and not username = ? order by timestamp asc`)
	if err != nil {
		log.Fatal(err)
		return messages
	}
	defer stmt.Close()

	rows, err = stmt.Query(t, c, u)
	if err != nil {
		log.Fatal(err)
		return messages
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var username string
		var timestamp int64
		var channel string
		var content string

		err = rows.Scan(&id, &username, &timestamp, &channel, &content)
		if err != nil {
			log.Fatal(err)
			return messages
		}

		message := models.Message{ID: id, Username: username, Timestamp: timestamp, Content: content}

		messages = append(messages, message)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return messages
}

// Inserts a new Message, converting the Partial Message into a Message beforehand.
func (w *DBWrapper) insertMessage(m models.ParialMessage) {
	tx, err := w.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into messages(id, username, timestamp, channel, content) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	message := m.GetMessage()

	_, err = stmt.Exec(message.ID, message.Username, message.Timestamp, message.Channel, message.Content)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got new message", message.ID)
}
