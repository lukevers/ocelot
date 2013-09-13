package main

import (
	"net"
)

func VerifyNetmask(netmask *net.IPNet, address string) bool {
	return netmask.Contains(net.ParseIP(address))
}