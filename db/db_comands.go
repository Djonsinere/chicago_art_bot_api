package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Base_init_db() {
	db, err := sql.Open("sqlite3", "./chicago_users.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB initialize")
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (id string not null, username text);

	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func CreateNewUser(id *string, username *string) {
	db, err := sql.Open("sqlite3", "./chicago_users.db")
	if err != nil {
		log.Fatal(err)
	}
	user_exist := is_user_exist(id, username)

	// Create
	if !user_exist {
		statement, _ := db.Prepare("INSERT INTO users (id, username) VALUES (?, ?)")
		statement.Exec(*id, *username)
		fmt.Println("Insert user into db", *id, *username)
	}
}

func is_user_exist(id *string, username *string) bool {
	db, error := sql.Open("sqlite3", "./chicago_users.db")
	if error != nil {
		log.Fatal(error)
	}
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND username = $2)`, *id, *username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}
