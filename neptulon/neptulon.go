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
	errMutex   sync.Mutex
	listener   *Listener
	middleware []func(conn *Conn, session *Session, msg []byte)
	conns      map[string]*Conn
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
	err := n.listener.Accept(handleConn(n), handleMsg(n), handleDisconn(n))
	if err != nil && n.debug {
		log.Fatalln("Listener returned an error while closing:", err)
	}

	n.errMutex.Lock()
	n.err = err
	n.errMutex.Unlock()

	return err
}

// Stop stops a server instance.
func (n *Neptulon) Stop() error {
	err := n.listener.Close()

	// close all active connections discarding any read/writes that is going on currently
	// this is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server
	for _, conn := range n.conns {
		err := conn.Close()
		if err != nil {
			return err
		}
	}

	n.errMutex.Lock()
	if n.err != nil {
		return n.err
	}
	n.errMutex.Unlock()
	return err
}

// handleConn handles client connected event.
func handleConn(n *Neptulon) func(conn *Conn, session *Session) {
	return func(conn *Conn, session *Session) {
		n.conns[session.id] = conn
	}
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
func handleDisconn(n *Neptulon) func(conn *Conn, session *Session) {
	return func(conn *Conn, session *Session) {
		delete(n.conns, session.id)
	}
}
