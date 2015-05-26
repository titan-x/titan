// Package neptulon is a socket framework with middleware support.
package neptulon

import (
	"log"
	"net/http"
	"sync"
)

// Neptulon framework entry point.
type Neptulon struct {
	debug       bool
	err         error
	listener    *Listener
	mutex       sync.Mutex
	middlewares []*func(ctx Context) (response interface{})
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
		mutex:    sync.Mutex{},
	}, nil
}

// Handle registers a new middleware to handle incoming messages.
func (n *Neptulon) Handle() {}

// Run starts accepting connections on the internal listener and handles connections with registered middleware.
// This function blocks and never returns, unless there was an error while accepting a new connection or the listner was closed.
func (n *Neptulon) Run() error {
	err := n.listener.Accept(handleMsg, handleDisconn)
	if err != nil && n.debug {
		log.Fatalln("Listener returned an error while closing:", err)
	}

	n.mutex.Lock()
	n.err = err
	n.mutex.Unlock()

	return err
}

// handleMsg handles incoming client messages.
func handleMsg(conn *Conn, session *Session, msg []byte) {

}

// handleDisconn handles client disconnection.
func handleDisconn(conn *Conn, session *Session) {}

// UseTLS enables TLS on a Neptulon app.
// func (n *Neptulon) UseTLS(cert, privKey []byte) {}

type handler func(w http.ResponseWriter, r *http.Request) error

func handle(handlers ...handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			err := handler(w, r)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}
	})
}
