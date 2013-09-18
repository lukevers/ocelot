package main

import (
	"database/sql"
)

var (
	Db DB
)

type DB struct {
	*sql.DB
	ReadOnly bool
}

func (db DB) InitalizeTables() (err error) {
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS users (
address BINARY(16) PRIMARY KEY,
name VARCHAR(255) NOT NULL,
description VARCHAR(255),
website VARCHAR(255));`)
	
	return
}