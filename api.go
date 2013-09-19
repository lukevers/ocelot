package main

import (
	"time"
	"math/rand"
	"net/http"
	"fmt"
)

var (
	ActiveTokens = make(map[uint32]token)
)

type token struct {
	IP string
	Issued time.Time
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