// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import "github.com/nbusy/devastator/neptulon"

// App is a Neptulon JSON-RPC app.
type App struct {
	middleware []func(conn *neptulon.Conn, session *neptulon.Session, msg *Message)
}

// NewApp creates a Neptulon JSON-RPC app.
func NewApp() {
}
