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