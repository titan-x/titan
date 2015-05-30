// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import (
	"encoding/json"

	"github.com/nbusy/devastator/neptulon"
)

// App is a Neptulon JSON-RPC app.
type App struct {
	middleware []func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)
}

// NewApp creates a Neptulon JSON-RPC app.
func NewApp(n *neptulon.App) *App {
	a := App{}
	n.Middleware(a.handler)
	return &a
}

// Middleware registers a new middleware to handle incoming messages.
func (a *App) Middleware(middleware func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)) {
	a.middleware = append(a.middleware, middleware)
}

func (a *App) handler(conn *neptulon.Conn, session *neptulon.Session, msg []byte) {
	var m Message
	if err := json.Unmarshal(msg, &m); err != nil {
		return
	}

	for _, mid := range a.middleware {
		mid(conn, session, &m)
	}
}