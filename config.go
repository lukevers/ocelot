package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	
	// Port is an integer that defines what port to run the
	// webserver for ocelot on
	Port int
	
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