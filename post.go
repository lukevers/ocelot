package main

import (
	"time"
)

type Post struct {
	ID          int
	Title       string
	User        *User
	Description string
	Posted      time.Time
	Body        string
}