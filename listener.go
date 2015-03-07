package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	ping  = []byte("ping")
	close = []byte("close")
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
		laddr = "localhost:443"
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
func (l *Listener) Accept(handleMsg func(msg []byte)) error {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			return fmt.Errorf("error while accepting a new connection from a client: %v", err)
			// todo: it might not be appropriate to break the loop on recoverable errors (like client disconnect during handshake)
			// the underlying fd.accept() does some basic recovery though we might need more: http://golang.org/src/net/fd_unix.go
		}

		log.Println("Client connected: listening for messages from client IP:", conn.RemoteAddr())
		go handleConn(conn, handleMsg)
	}
}

func handleConn(conn net.Conn, handleMsg func(msg []byte)) {
	defer conn.Close()
	defer log.Println("Closed connection to client with IP:", conn.RemoteAddr())
	buf := make([]byte, 4096) // same limit as Google Cloud Messaging for simplicity
	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Minute * 5))
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			log.Fatalln("Client read error: ", err)
			break
		}
		log.Printf("Read %v bytes from client with IP: %v\n", n, conn.RemoteAddr())

		if n == 4 && bytes.Equal(buf[:n], ping) {
			continue
		} else if n == 5 && bytes.Equal(buf[:n], close) {
			go handleMsg(buf[:n])
			return
		}

		go handleMsg(buf[:n])
	}
}

// Close closes the listener.
func (l *Listener) Close() error {
	defer log.Println("Listener was closed on local network address:", l.listener.Addr())
	return l.listener.Close()
}
