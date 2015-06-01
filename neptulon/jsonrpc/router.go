package jsonrpc

import (
	"log"

	"github.com/nbusy/devastator/neptulon"
)

// Router is a JSON-RPC request routing middleware.
type Router struct {
	routes map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg *Message) (result interface{}, err *ResError)
}

// NewRouter creates a JSON-RPC router instance and registers it with the Neptulon JSON-RPC app.
func NewRouter(app *App) (*Router, error) {
	r := Router{
		routes: make(map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg *Message) (result interface{}, err *ResError)),
	}

	app.Middleware(r.middleware)

	return &r, nil
}

// Route adds a new route registry.
func (r *Router) Route(route string, handler func(conn *neptulon.Conn, session *neptulon.Session, msg *Message) (result interface{}, err *ResError)) {
	r.routes[route] = handler
}

func (r *Router) middleware(conn *neptulon.Conn, session *neptulon.Session, msg *Message) {
	if res, err := r.routes[msg.Method](conn, session, msg); res != nil && msg.ID != "" {
		if n, err := conn.WriteMsg(Response{ID: msg.ID, Result: res}); err != nil {
			log.Fatalln("Errored while sending JSON-RPC response:", err)
		} else if n == 0 {
			log.Fatalln("Errored while sending JSON-RPC response: wrote 0 bytes to connection")
		}
	} else if err != nil && msg.ID != "" {
		if n, e := conn.WriteMsg(Response{ID: msg.ID, Error: err}); e != nil {
			log.Fatalln("Errored while sending JSON-RPC error response:", err)
		} else if n == 0 {
			log.Fatalln("Errored while sending JSON-RPC error response: wrote 0 bytes to connection")
		}
	}
}
