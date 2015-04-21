package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var (
	ping   = []byte("ping")
	closed = []byte("close")
)

// Listener accepts connections from devices.
type Listener struct {
	debug    bool
	listener net.Listener
	wg       sync.WaitGroup
}

// Listen creates a TCP listener with the given PEM encoded X.509 certificate and the private key on the local network address laddr.
// Debug mode logs all server activity.
func Listen(cert, privKey []byte, laddr string, debug bool) (*Listener, error) {
	tlsCert, err := tls.X509KeyPair(cert, privKey)
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cert)
	if err != nil || !ok {
		return nil, fmt.Errorf("failed to parse the certificate or the private key: %v", err)
	}

	conf := tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientCAs:    pool,
		ClientAuth:   tls.VerifyClientCertIfGiven,
	}

	l, err := tls.Listen("tcp", laddr, &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS listener no network address %v: %v", laddr, err)
	}
	if debug {
		log.Printf("Listener created with local network address: %v\n", laddr)
	}

	return &Listener{
		debug:    debug,
		listener: l,
	}, nil
}

// Session is a generic session data store for client handlers.
type Session struct {
	UserID       uint32
	Error        string
	Data         interface{}
	Disconnected bool
}

// Accept waits for incoming connections and forwards the client connect/message/disconnect events to provided handlers in a new goroutine.
// This function blocks and never returns, unless there is an error while accepting a new connection.
func (l *Listener) Accept(handleMsg func(conn *Conn, session *Session, msg []byte), handleDisconn func(conn *Conn, session *Session)) error {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			return fmt.Errorf("error while accepting a new connection from a client: %v", err)
			// todo: it might not be appropriate to break the loop on recoverable errors (like client disconnect during handshake)
			// the underlying fd.accept() does some basic recovery though we might need more: http://golang.org/src/net/fd_unix.go
		}
		tlsconn, ok := conn.(*tls.Conn)
		if !ok {
			return errors.New("cannot cast net.Conn interface to tls.Conn type")
		}
		if l.debug {
			log.Println("Client connected: listening for messages from client IP:", conn.RemoteAddr())
		}
		go handleClient(&l.wg, NewConn(tlsconn, 0, 0), l.debug, handleMsg, handleDisconn)
	}
}

// handleClient waits for messages from the connected client and forwards the client message/disconnect
// events to provided handlers in a new goroutine.
// This function never returns, unless there is an error while reading from the channel or the client disconnects.
func handleClient(wg *sync.WaitGroup, conn *Conn, debug bool, handleMsg func(conn *Conn, session *Session, msg []byte), handleDisconn func(conn *Conn, session *Session)) {
	wg.Add(1)
	defer wg.Done()

	session := &Session{}

	if debug {
		defer func() {
			if session.Disconnected {
				log.Println("Client disconnected on IP:", conn.RemoteAddr())
			} else {
				log.Println("Closed connection to client with IP:", conn.RemoteAddr())
			}
		}()
	}
	defer conn.Close() // todo: handle close error, store the error in conn object and return it to handleMsg/handleErr/handleDisconn or one level up (to server)

	for {
		if session.Error != "" {
			// todo: send error message to user, log the error, and close the conn and return
			return
		}

		n, msg, err := conn.Read()
		if err != nil {
			if err == io.EOF {
				session.Disconnected = true
				break
			}

			log.Fatal("errored while reading:", err)
		}

		// shortcut 'ping' and 'close' messages, saves some processing time
		if n == 4 && bytes.Equal(msg, ping) {
			continue // send back pong?
		}
		if n == 5 && bytes.Equal(msg, closed) {
			go func() {
				wg.Add(1)
				defer wg.Done()
				handleDisconn(conn, session)
			}()
			return
		}

		go func() {
			wg.Add(1)
			defer wg.Done()
			handleMsg(conn, session, msg)
		}()
	}
}

// Close closes the listener.
func (l *Listener) Close() error {
	l.wg.Wait()
	if l.debug {
		defer log.Println("Listener was closed on local network address:", l.listener.Addr())
	}
	return l.listener.Close()
}
