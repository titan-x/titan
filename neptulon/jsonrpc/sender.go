package jsonrpc

import "github.com/nbusy/devastator/neptulon"

// Sender is a JSON-RPC request sending middleware. This should be registered before any incoming request routers.
type Sender struct {
	routes map[string]func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)
}

func (s *Sender) middleware(conn *neptulon.Conn, session *neptulon.Session, msg *Message) {
	s.routes[msg.Method](conn, session, msg)
}
