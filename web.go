package main

import (
	"net/http"
)

func Serve(config *Config) {	

	http.HandleFunc("/", HandleRoot)
	
	l.Infof("Starting server on %s", config.Address)
	err := http.ListenAndServe(config.Address, nil)
	if err != nil {
		l.Fatalf("Listen and serve error: %s", err)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	
}