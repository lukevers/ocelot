package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	
	// Address contains the interface, and port all in one to save
	// time instead of combinding multiple variables into one. For
	// example: interface:port
	Address string
	
	// Netmask is for determining if the incoming address is the
	// same type of netmask that it should be in order to sign in
	// and sign up.
	Netmask string
	
	// The Database struct contains information about the
	// database.
	Database struct {
		// ReadOnly is a bool that is either true or false. If
		// it is rtrue then we don't allow any editing to the
		// database.
		ReadOnly bool
		
		// Driver is the type of database that we are
		// using. For example, by default we use sqlite3.
		Driver string
		
		// Resource is the name of the database that we are
		// looking for to use. 
		Resource string
	}
}

// ReadConfig reads the configuration file from JSON and returns it in
// the form of a *Config.
func ReadConfig(path string) (config *Config, err error) {
	file, err := os.Open(path)
	defer file.Close()
	
	if err != nil {
		return
	}
	
	config = &Config{}
	err = json.NewDecoder(file).Decode(config)
	
	return 
}