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

// HandleRoot is the default handler that figures out if you're
// already signed up or not. If you're not on Hyperboria, you can't do
// anything. If you're on Hyperboria and you're not signed up yet,
// it'll give you the signup page. If you're on Hyperboria and you're
// already signed up, it'll automatically sign you in and show you the
// front page.
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	r.RemoteAddr, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		l.Noticef("SplitHostPort error: %s", err)
	}
	
	templates = template.Must(template.ParseGlob("templates/*"))

	// If we are not on Hyperboria, send them to the "nohype"
	// page. If we are, check if we are a user yet.
	if VerifyNetmask(Netmask, r.RemoteAddr) {
		if VerifiedUser(r.RemoteAddr) {
			templates.ExecuteTemplate(w, "index", nil)
		} else {
			// Create a temporary struct to put the
			// address into the form to create an account.
			type Address struct {
				Address string
			}
			templates.ExecuteTemplate(w, "signup", 
				&Address{Address: r.RemoteAddr})
		}
	} else {
		// Show a blank page for non-hype right now.
		templates.ExecuteTemplate(w, "nohype", nil)
	}
}

// HandleAssets is a static file server that serves everything in the
// assets directory.
func HandleAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}