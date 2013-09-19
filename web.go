package main

import (
	"net/http"
	"net"
	"html/template"
	"fmt"
	"strconv"
)

var (
	Netmask *net.IPNet
	err error
	templates *template.Template
)

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
	http.HandleFunc("/settings", HandleSettings)
	http.HandleFunc("/me", HandleMe)
	http.HandleFunc("/create", HandleCreate)

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

// HandleRoot is the handler for the general / page
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	r.RemoteAddr, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		l.Noticef("SplitHostPort error: %s", err)
	}
	
	templates = template.Must(template.ParseGlob("templates/*"))

	// If we are not on Hyperboria, send them to the "nohype"
	// page. If we are, check if we are a user yet.
	if VerifyNetmask(Netmask, r.RemoteAddr) {
		exists, user := VerifiedUser(r.RemoteAddr)
		if exists {
			templates.ExecuteTemplate(w, "index", &user)
		} else {
			// Redirect to the signup page
			SignUp(w, r.RemoteAddr, templates)
		}
	} else {
		// Redirect to the no access page
		NoAccess(w, templates)
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
	} else if fname == "" || fname == " " {
		fmt.Fprintf(w, "Error: First name can not be empty")
	} else if lname == "" || lname == " " {
		fmt.Fprintf(w, "Error: Last name can not be empty")
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

// HandleSettings is the handler for the /settings page
func HandleSettings(w http.ResponseWriter, r *http.Request) {
	r.RemoteAddr, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		l.Noticef("SplitHostPort error: %s", err)
	}
	
	templates = template.Must(template.ParseGlob("templates/*"))
	
	// If we are not signed in and on Hyperboria send to the
	// signup page. If we are not signed in and not on Hyperboria
	// send to the nohype page.
	if VerifyNetmask(Netmask, r.RemoteAddr) {
		exists, user := VerifiedUser(r.RemoteAddr)
		if exists {
			templates.ExecuteTemplate(w, "settings", &user)
		} else {
			// Redirect to signup page
			SignUp(w, r.RemoteAddr, templates)
		}
	} else {
		// Redirect to no access page
		NoAccess(w, templates)
	}
}

// HandleMe is the handler for the /me page
func HandleMe(w http.ResponseWriter, r *http.Request) {

}

// HandleCreate is the handler for the /create page
func HandleCreate(w http.ResponseWriter, r *http.Request) {

}

// SignUp is a redirect to the signup page
func SignUp(w http.ResponseWriter, address string, templates *template.Template) {
	// Create a temporary struct to put the address into the form
	// to create an account.
	type Address struct {
		Address string
	}
	templates.ExecuteTemplate(w, "signup", &Address{Address: address})
}

// NoAccess is a redirect to the noaccess page
func NoAccess(w http.ResponseWriter, templates *template.Template) {
	templates.ExecuteTemplate(w, "nohype", nil)
}