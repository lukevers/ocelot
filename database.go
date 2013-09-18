package main

import (
	"database/sql"
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"time"
)

var (
	Db DB
)

type DB struct {
	*sql.DB
	ReadOnly bool
}

// InitalizeTables creates the tables needed in the database if they
// do not already exist.
func (db DB) InitalizeTables() (err error) {
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS users (
address BINARY(16) PRIMARY KEY,
name VARCHAR(255) NOT NULL,
description VARCHAR(255),
website VARCHAR(255),
updated INT NOT NULL);`)
	
	return
}

// LengthUsers counts all of the users in the database. If there is an
// error it returns -1.
func (db DB) LengthUsers() (length int) {
	row := db.QueryRow("SELECT COUNT(*) FROM users;")
	if err := row.Scan(&length); err != nil {
		l.Errf("Error counting number of users: %s", err)
		length = -1
	}
	return
}

// AddUser adds a new user to the database. If there is an error it
// returns the error.
func (db DB) AddUser(user *User) (err error) {
	// Prepare database
	u, err := db.Prepare(`INSERT INTO users
(address, name, description, website, updated)
VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		return
	}
	// Insert values
	_, err = u.Exec([]byte(user.Address), user.Name, user.Description,
		user.Website, time.Now())
	defer u.Close()
	return
}
