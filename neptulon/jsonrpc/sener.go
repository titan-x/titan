package jsonrcp

import "github.com/nbusy/devastator/neptulon"

// Sender is a JSON-RPC request sending middleware. This should be registered before any incoming request routers.
type Sender struct {
	routes map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)
}

func (r *Router) middleware(conn *neptulon.Conn, session *neptulon.Session, msg *Message) {
	r.routes[msg.Method](conn, session, msg)
}
