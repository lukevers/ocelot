package main

import (
	"net/http"
	"net"
	"html/template"
	"fmt"
	"time"
	"math/rand"
	"strconv"
)

var (
	Netmask *net.IPNet
	err error
	templates *template.Template
	ActiveTokens = make(map[uint32]token)
)

type token struct {
	IP string
	Issued time.Time
}

// Serve takes the options from the configuration file and starts up
// the servers.
func Serve(config *Config) {	
	// Check Netmask
	_, Netmask, err = net.ParseCIDR(config.Netmask)
	if err != nil {
		l.Fatalf("Could not parse netmask error: %s", err)
	}
	
	// Handle normal pages
	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/assets/", HandleAssets)

	// Handle API points
	http.HandleFunc("/api/newuser", HandleNewUser)
	http.HandleFunc("/api/token", HandleToken)

	// Start server
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

// HandleNewUser is a handler that creates a new user when there is a
// POST to /api/newuser.
func HandleNewUser(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	u := &User{
		Address: r.FormValue("address"),
		Name: CreateFullName(fname, lname),
		Description: r.FormValue("desc"),
		Website: r.FormValue("website"),
	}
	t, err := strconv.ParseUint(r.FormValue("token"), 10, 0)
	if err != nil {
		l.Emergf("Could not parse token to uint: %s", err)
	}
	
	if !CheckToken(r.FormValue("address"), uint32(t)) {
		fmt.Fprintf(w, "Error: Token invalid")
	} else if !VerifyNetmask(Netmask, u.Address) {
		fmt.Fprintf(w, "Error: Address does not match netmask")
	} else if u.Name == " " {
		fmt.Fprintf(w, "Error: Name can not be empty")
	} else if len(u.Name) > 255 {
		fmt.Fprintf(w, "Error: Name can not be longer than 255 characters")
	} else if len(u.Description) > 255 {
		fmt.Fprintf(w, "Error: Description can not be longer than 255 characters")
	} else if len(u.Website) > 255 {
		fmt.Fprintf(w, "Error: Website can not be longer than 255 characters")
	} else {
		err := Db.AddUser(u)
		if err != nil {
			l.Emergf("Could not add user to database: %s", err)
			fmt.Fprintf(w, "Error: %s", GetDetailedError(err))
		} else {
			fmt.Fprintf(w, "Success: User added")
			l.Infof("User %s added to database", u.Name)
		}
	}
}

// HandleToken generates a random token and keeps it in memory.
func HandleToken(w http.ResponseWriter, r *http.Request) {
	tokenid := rand.Uint32()
	ActiveTokens[tokenid] = token{r.FormValue("address"), time.Now()}
	fmt.Fprint(w, tokenid)
}

// CheckToken ensures that a particular token is valid, meaning that
// it is in the list, and has not yet expired. If so, it removes the
// token and returns true. If anything else occurs, it returns
// false. Tokens are active for 1 minute after they are created.
func CheckToken(IP string, token uint32) bool {
	t, ok := ActiveTokens[token]
	if ok {
		delete(ActiveTokens, token)
	}
	
	if !ok || time.Now().After(t.Issued.Add(time.Minute)) ||
		t.IP != IP {
		return false
	}

	return true
}