package database

import (
	"database/sql"
	"fmt"
	"game-ranker/users-manager/internal"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type User = internal.User

// Will only create the DB table if it doesn't exist
func InitDbTable() {
	fmt.Println("Initiating table")
	db := Connect()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		public_id TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		pass TEXT NOT NULL,
		email TEXT NOT NULL
	);
	`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	defer Close(db)
}

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Close(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func AddUser(user User) error {
	db := Connect()
	defer Close(db)

	sqlStmt := `
		INSERT INTO users (public_id, username, pass, email)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(sqlStmt, user.ID, user.Username, user.HashedPass, user.Email)
	if err != nil {
		return err
	}

	return nil
}
