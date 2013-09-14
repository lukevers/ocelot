package main

import (
	"net/http"
	"net"
	"fmt"
)

var (
	Netmask *net.IPNet
	err error
)

// Serve takes the options from the configuration file and starts up
// the servers.
func Serve(config *Config) {	

	_, Netmask, err = net.ParseCIDR(config.Netmask)
	if err != nil {
		l.Fatalf("Could not parse netmask error: %s", err)
	}

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/assets/", HandleAssets)

	l.Infof("Starting server on %s", config.Address)
	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		l.Fatalf("Listen and serve error: %s", err)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	addr, hype, err := VerifyNetmask(Netmask, r.RemoteAddr)
	r.RemoteAddr = addr
	if err != nil {
		l.Noticef("Verify netmask error: %s", err)
	}

	// Now do hype/non-hype related things!
	if hype {
		fmt.Fprintf(w, "Hello, %q, you are on hype!", r.RemoteAddr)
	} else {
		fmt.Fprint(w, "You are not on hype!")
	}
}

// HandleAssets is a static file server that serves everything in the
// assets directory.
func HandleAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}