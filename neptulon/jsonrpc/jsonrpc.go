// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import (
	"encoding/json"

	"github.com/nbusy/devastator/neptulon"
)

// App is a Neptulon JSON-RPC app.
type App struct {
	middleware []func(conn *neptulon.Conn, msg *Message)
}

// NewApp creates a Neptulon JSON-RPC app.
func NewApp(n *neptulon.App) (*App, error) {
	a := App{}
	n.Middleware(a.handler)
	return &a, nil
}

// Middleware registers a new middleware to handle incoming messages.
func (a *App) Middleware(middleware func(conn *neptulon.Conn, msg *Message)) {
	a.middleware = append(a.middleware, middleware)
}

func (a *App) handler(conn *neptulon.Conn, msg []byte) {
	var m Message
	if err := json.Unmarshal(msg, &m); err != nil {
		return
	}

	// todo: it might be better to use Koa like stack: mid(conn, &m, mid(conn, &m, conn(....)))
	for _, mid := range a.middleware {
		mid(conn, &m)
	}
}
