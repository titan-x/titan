// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import (
	"encoding/json"
	"log"

	"github.com/nbusy/devastator/neptulon"
)

// App is a Neptulon JSON-RPC app.
type App struct {
	middleware []func(conn *neptulon.Conn, msg *Message) (result interface{}, err *ResError)
}

// NewApp creates a Neptulon JSON-RPC app.
func NewApp(n *neptulon.App) (*App, error) {
	a := App{}
	n.Middleware(a.handler)
	return &a, nil
}

// Middleware registers a new middleware to handle incoming messages.
func (a *App) Middleware(middleware func(conn *neptulon.Conn, msg *Message) (result interface{}, err *ResError)) {
	a.middleware = append(a.middleware, middleware)
}

func (a *App) handler(conn *neptulon.Conn, msg []byte) {
	var m Message
	if err := json.Unmarshal(msg, &m); err != nil {
		log.Fatalln("Cannot deserialize incoming message:", err)
	}

	for _, mid := range a.middleware {
		if res, err := mid(conn, &m); res != nil {
			if m.Method == "" || m.ID == "" {
				log.Fatalln("Cannot return a response to a non request")
			}

			if _, err := conn.WriteMsg(Response{ID: m.ID, Result: res, Error: err}); err != nil {
				log.Fatalln("Errored while sending JSON-RPC response:", err)
			}
			break // todo: evaluate using Koa like stack instead of manual break: mid(conn, &m, mid(conn, &m, conn(....)))
		}
	}
}
