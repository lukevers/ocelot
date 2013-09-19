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
	if err != nil {
		return
	}

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS posts (
title VARCHAR(255) PRIMARY KEY,
address BINARY(16) NOT NULL,
description VARCHAR(255) NOT NULL,
posted INT NOT NULL);`)
	
	return
}

// Lengh counts all of the `x` in the database. Depending on what you
// want to count. You call the Length func by giving the table name as
// a parameter.
func (db DB) Length(table string) (length int) {
	row := db.QueryRow("SELECT COUNT(*) FROM "+table+";")
	if err := row.Scan(&length); err != nil {
		l.Errf("Error counting number of %s: %s", table, err)
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

// GetUser checks the users table in the database for a user. If the
// user is there, it returns a *User and and a nil error. It will
// return an actual error and a nil *User if there are errors.
func (db DB) GetUser(address string) (user *User, err error) {
	// Prepare database
	u, err := db.Prepare(`
SELECT name, description, website
FROM users
WHERE address = ?
LIMIT 1`)
	if err != nil {
		return
	}
	
	user = &User{Address: address}
	baddr := []byte(address)
	
	// Query database
	row := u.QueryRow(baddr)
	err = row.Scan(&user.Name, &user.Description, &user.Website)
	defer u.Close()
	
	// If the error is of the particular type sql.ErrNoRows, it
	// simply means that the node does not exist. In that case,
	// return (nil, nil).
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return
}