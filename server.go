package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
)

// Listener accepts connections from devices.
type Listener struct {
	debug    bool
	listener *net.Listener
}

// Listen creates a listener with the given PEM encoded X.509 certificate and the private key on the local network address laddr.
// Debug mode logs all server activity.
func Listen(cert, priv []byte, laddr string, debug bool) (*Listener, error) {
	c, err := x509.ParseCertificate(cert)
	p, err := x509.ParsePKCS1PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the certificate or the private key with error: %v", err)
	}

	pool := x509.NewCertPool()
	pool.AddCert(c)

	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  p,
	}

	config := tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientCAs:    pool,
	}

	if laddr == "" {
		laddr = "0.0.0.0:443"
	}

	listener, err := tls.Listen("tcp", laddr, &config)
	if err != nil {
		return nil, err
	}

	if debug {
		log.Printf("New CCS connection established with XMPP parameters: %+v\n", c)
	}

	return &Listener{
		debug:    debug,
		listener: &listener,
	}, nil
}
