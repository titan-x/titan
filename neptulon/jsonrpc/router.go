package jsonrpc

import (
	"encoding/json"

	"github.com/nbusy/devastator/neptulon"
)

// NewRouter creates a JSON-RPC 2.0 router instance and registers it with the Neptulon app.
func NewRouter(app *neptulon.App) *Router {
	r := Router{}
	app.Middleware(r.middleware)
	return &r
}

// Router is a JSON-RPC request routing middleware.
type Router struct {
	routes map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)
}

// Route adds a new route registry.
func (r *Router) Route(route string, handler func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)) {
	r.routes[route] = handler
}

func (r *Router) middleware(conn *neptulon.Conn, session *neptulon.Session, msg []byte) {
	var m Message
	if err := json.Unmarshal(msg, &m); err != nil {
		return
	}

	r.routes[m.Method](conn, session, &m)
}
