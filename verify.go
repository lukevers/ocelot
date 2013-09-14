package main

import (
	"net"
)

func VerifyNetmask(netmask *net.IPNet, address string) (string, bool, error) {
	// Get rid of [host]:port form and only keep the host. We
	// don't need to know the port.
	address, _, err = net.SplitHostPort(address)
	if err != nil {
		return address, false, err
	}
	return address, netmask.Contains(net.ParseIP(address)), nil
}