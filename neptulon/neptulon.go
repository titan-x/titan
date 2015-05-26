// Package neptulon is a socket framework with middleware support.
package neptulon

import (
	"log"
	"sync"
)

// Neptulon framework entry point.
type Neptulon struct {
	debug      bool
	err        error
	listener   *Listener
	errMutex   sync.RWMutex
	middleware []func(conn *Conn, session *Session, msg []byte)
}

// New creates and returns a new Neptulon app. This is the default TLS constructor.
// Debug mode dumps raw TCP data to stderr (log.Println() default).
func New(cert, privKey []byte, laddr string, debug bool) (*Neptulon, error) {
	l, err := Listen(cert, privKey, laddr, debug)
	if err != nil {
		return nil, err
	}

	return &Neptulon{
		debug:    debug,
		listener: l,
	}, nil
}

// Middleware registers a new middleware to handle incoming messages.
func (n *Neptulon) Middleware(middleware func(conn *Conn, session *Session, msg []byte)) {
	n.middleware = append(n.middleware, middleware)
}

// Run starts accepting connections on the internal listener and handles connections with registered middleware.
// This function blocks and never returns, unless there was an error while accepting a new connection or the listner was closed.
func (n *Neptulon) Run() error {
	err := n.listener.Accept(handleMsg(n), handleDisconn)
	if err != nil && n.debug {
		log.Fatalln("Listener returned an error while closing:", err)
	}

	n.errMutex.Lock()
	n.err = err
	n.errMutex.Unlock()

	return err
}

// Stop stops a server instance.
func (s *Server) Stop() error {
	err := s.listener.Close()

	// close all active connections discarding any read/writes that is going on currently
	// this is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server
	for _, user := range users {
		err := user.Conn.Close()
		if err != nil {
			return err
		}
		user.Conn = nil
	}
	for _, conn := range s.listener.Conns {
		err := conn.Close()
		if err != nil {
			return err
		}
	}

	s.mutex.Lock()
	if s.err != nil {
		return s.err
	}
	s.mutex.Unlock()
	return err
}

// handleMsg handles incoming client messages.
func handleMsg(n *Neptulon) func(conn *Conn, session *Session, msg []byte) {
	return func(conn *Conn, session *Session, msg []byte) {
		for _, m := range n.middleware {
			m(conn, session, msg)
		}
	}
}

// handleDisconn handles client disconnection.
func handleDisconn(conn *Conn, session *Session) {}
