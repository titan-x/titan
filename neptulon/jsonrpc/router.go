package jsonrpc

import "github.com/nbusy/devastator/neptulon"

// register routes with callbacks here

// handle anonymouse calls here

// handle authentication here

// handle authenticated calls here

// Router is a simple routing middleware.
type Router struct {
	routes map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg interface{})
}

// Register adds a new route registry.
func (r *Router) Register(route string, fn func(conn *neptulon.Conn, session *neptulon.Session, msg interface{})) {
	r.routes[route] = fn
}
