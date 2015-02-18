package main

import (
	"crypto/tls"
	"net"
)

// Listener accepts connections from devices.
type Listener struct {
	tls *net.Listener
}

// Listen creates a listener with given server certificate. Debug mode logs all communications.
func Listen(cert tls.Certificate, debug bool) (*Listener, error) {
	return nil, nil
}
