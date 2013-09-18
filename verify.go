package main

import (
	"net"
)

// VerifyNetmask returns true if the address a user is coming from
// matches the netmask defined in the configuration file, passed to
// VerifyNetmask as `netmask *net.IPNet`
func VerifyNetmask(netmask *net.IPNet, address string) bool {
	return netmask.Contains(net.ParseIP(address))
}

// VerifiedUser checks if the user exists in the database. If the user
// exists then we return true, else false. VerifyNetmask will always
// be called before we run this function, so we don't have to check
// the netmask again before checking if the user exists.
func VerifiedUser(address string) (bool, *User) {	
	user, err := Db.GetUser(address)
	if err != nil {
		return false, nil
	}
	
	if user != nil {
		return true, user
	}
	
	return false, nil
}