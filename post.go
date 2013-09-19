package main

import (
	"time"
)

type Post struct {
	Title       string
	User        *User
	Description string
	Posted      time.Time
	Body        string
}