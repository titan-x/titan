// Package neptulon is a socket framework with middleware support.
package neptulon

import (
	"log"
	"sync"
)

// App is a Neptulon application.
type App struct {
	debug      bool
	err        error
	errMutex   sync.RWMutex
	listener   *Listener
	middleware []func(conn *Conn, session *Session, msg []byte)
	conns      map[string]*Conn
}

// NewApp creates and returns a new Neptulon app. This is the default TLS constructor.
// Debug mode dumps raw TCP data to stderr (log.Println() default).
func NewApp(cert, privKey []byte, laddr string, debug bool) (*App, error) {
	l, err := Listen(cert, privKey, laddr, debug)
	if err != nil {
		return nil, err
	}

	return &App{
		debug:    debug,
		listener: l,
		conns:    make(map[string]*Conn),
	}, nil
}

// Middleware registers a new middleware to handle incoming messages.
func (a *App) Middleware(middleware func(conn *Conn, session *Session, msg []byte)) {
	a.middleware = append(a.middleware, middleware)
}

// Run starts accepting connections on the internal listener and handles connections with registered middleware.
// This function blocks and never returns, unless there was an error while accepting a new connection or the listner was closed.
func (a *App) Run() error {
	err := a.listener.Accept(handleConn(a), handleMsg(a), handleDisconn(a))
	if err != nil && a.debug {
		log.Fatalln("Listener returned an error while closing:", err)
	}

	a.errMutex.Lock()
	a.err = err
	a.errMutex.Unlock()

	return err
}

// Stop stops a server instance.
func (a *App) Stop() error {
	err := a.listener.Close()

	// close all active connections discarding any read/writes that is going on currently
	// this is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server
	for _, conn := range a.conns {
		err := conn.Close()
		if err != nil {
			return err
		}
	}

	a.errMutex.RLock()
	if a.err != nil {
		return a.err
	}
	a.errMutex.RUnlock()
	return err
}

func handleConn(a *App) func(conn *Conn, session *Session) {
	return func(conn *Conn, session *Session) {
		a.conns[session.id] = conn
	}
}

func handleMsg(a *App) func(conn *Conn, session *Session, msg []byte) {
	return func(conn *Conn, session *Session, msg []byte) {
		for _, m := range a.middleware {
			m(conn, session, msg)
		}
	}
}

func handleDisconn(a *App) func(conn *Conn, session *Session) {
	return func(conn *Conn, session *Session) {
		delete(a.conns, session.id)
	}
}
