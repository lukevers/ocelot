package main

import (
	"net/http"
	"net"
	"html/template"
)

var (
	Netmask *net.IPNet
	err error
	templates *template.Template
)

type Page struct {
	User *User
}

type User struct {
	IPv6 string
	Name string
}

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
	r.RemoteAddr, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		l.Noticef("SplitHostPort error: %s", err)
	}
	
	templates = template.Must(template.ParseGlob("templates/*"))
	
	

	// Now do hype/non-hype related things!
	if VerifyNetmask(Netmask, r.RemoteAddr) {
		templates.ExecuteTemplate(w, "index", nil)
	} else {
		templates.ExecuteTemplate(w, "index", nil)
	}
}

// HandleAssets is a static file server that serves everything in the
// assets directory.
func HandleAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}