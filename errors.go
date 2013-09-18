package main

import (
	"strings"
)

func GetDetailedError(err error) string {
	if strings.Contains(err.Error(), "[2067]") {
		return "Address is already taken"
	}
	
	// If we have yet to come across this error, send the regular
	// error and we'll eventually add it here.
	return err.Error()
}