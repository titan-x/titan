package jsonrpc

import "github.com/nbusy/devastator/neptulon"

// Router is a JSON-RPC request routing middleware.
type Router struct {
	routes map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)
}

// NewRouter creates a JSON-RPC router instance and registers it with the Neptulon JSON-RPC app.
func NewRouter(app *App) (*Router, error) {
	r := Router{
		routes: make(map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)),
	}

	app.Middleware(r.middleware)

	return &r, nil
}

// Route adds a new route registry.
func (r *Router) Route(route string, handler func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)) {
	r.routes[route] = handler
}

func (r *Router) middleware(conn *neptulon.Conn, session *neptulon.Session, msg *Message) {
	r.routes[msg.Method](conn, session, msg)
}
