package jsonrpc

import "github.com/nbusy/devastator/neptulon"

// register routes with callbacks here

// handle anonymouse calls here

// handle authentication here

// handle authenticated calls here

// Router is a JSON-RPC request routing middleware.
type Router struct {
	requestRoutes      map[string]func(conn *neptulon.Conn, session *neptulon.Session, req *Request)
	notificationRoutes map[string]func(conn *neptulon.Conn, session *neptulon.Session, not *Notification)
}

// RegisterReq adds a new request route registry.
func (r *Router) RegisterReq(route string, handler func(conn *neptulon.Conn, session *neptulon.Session, req *Request)) {
	r.requestRoutes[route] = handler
}

// RegisterNot adds a new request route registry.
func (r *Router) RegisterNot(route string, handler func(conn *neptulon.Conn, session *neptulon.Session, not *Notification)) {
	r.notificationRoutes[route] = handler
}
