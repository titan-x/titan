package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
)

// Listener accepts connections from devices.
type Listener struct {
	debug    bool
	listener net.Listener
}

// Listen creates a TCP listener with the given PEM encoded X.509 certificate and the private key on the local network address laddr.
// Debug mode logs all server activity.
func Listen(cert, privKey []byte, laddr string, debug bool) (*Listener, error) {
	tlsCert, err := tls.X509KeyPair(cert, privKey)
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cert)
	if err != nil || !ok {
		return nil, fmt.Errorf("failed to parse the certificate or the private key with error: %v", err)
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
		log.Printf("Listener created with local network address: %v\n", laddr)
	}

	return &Listener{
		debug:    debug,
		listener: listener,
	}, nil
}

// Accept waits for incoming connections and forwards incoming messages to handleMsg in a new goroutine.
// This function never returns, unless there is an error while accepting a new connection.
func (l *Listener) Accept( /*handleMsg func(string)*/ ) error {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			return fmt.Errorf("error while accepting a new connection from a client: %v", err)
			// todo: it might not be appropriate to break the loop on recoverable errors (like client disconnect during handshake)
			// the underlying fd.accept() does some basic recovery though we might need more: http://golang.org/src/net/fd_unix.go
		}

		log.Printf("Accepted connection and waiting for data from client IP: %v\n", conn.RemoteAddr())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	defer log.Printf("Closed connection to client with IP: %v\n", conn.RemoteAddr())
	for {
		written, err := io.Copy(conn, conn)
		if err != nil {
			break
		}
		log.Printf("Echoed %v bytes to client with IP: %v\n", written, conn.RemoteAddr())
	}
}

// Close closes the listener.
func (l *Listener) Close() error {
	defer log.Printf("Listener on local network address %v was closed.\n", l.listener.Addr())
	return l.listener.Close()
}
