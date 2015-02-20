package main

import (
	"crypto/x509"
	"net"
)

// Listener accepts connections from devices.
type Listener struct {
	tls *net.Listener
}

// Listen creates a listener with the given PEM encoded X.509 certificate and the private key. Debug mode logs all communications.
func Listen(cert, priv []byte, debug bool) (*Listener, error) {
	c, err := x509.ParseCertificate(cert)
	p, err := x509.ParsePKCS1PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
