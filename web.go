package main

import (
	"net/http"
)

func Serve(config *Config) {	

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/assets/", HandleAssets)

	l.Infof("Starting server on %s", config.Address)
	err := http.ListenAndServe(config.Address, nil)
	if err != nil {
		l.Fatalf("Listen and serve error: %s", err)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	
}

// HandleAssets is a static file server that serves everything in the
// assets directory.
func HandleAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}