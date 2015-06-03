package jsonrpc

import (
	"log"

	"github.com/nbusy/devastator/neptulon"
)

// Router is a JSON-RPC request routing middleware.
type Router struct {
	requestRoutes      map[string]func(conn *neptulon.Conn, req *Request) (result interface{}, err *ResError)
	notificationRoutes map[string]func(conn *neptulon.Conn, not *Notification)
}

// NewRouter creates a JSON-RPC router instance and registers it with the Neptulon JSON-RPC app.
func NewRouter(app *App) (*Router, error) {
	r := Router{
		requestRoutes:      make(map[string]func(conn *neptulon.Conn, req *Request) (result interface{}, err *ResError)),
		notificationRoutes: make(map[string]func(conn *neptulon.Conn, not *Notification)),
	}

	app.Middleware(r.middleware)
	return &r, nil
}

// Request adds a new request route registry.
func (r *Router) Request(route string, handler func(conn *neptulon.Conn, req *Request) (result interface{}, err *ResError)) {
	r.requestRoutes[route] = handler
}

// Notification adds a new notification route registry.
func (r *Router) Notification(route string, handler func(conn *neptulon.Conn, not *Notification)) {
	r.notificationRoutes[route] = handler
}

func (r *Router) middleware(conn *neptulon.Conn, msg *Message) {
	// if not request or notification
	if msg.Method == "" {
		return
	}

	// if request
	if msg.ID != "" {
		if handler, ok := r.requestRoutes[msg.Method]; ok {
			if res, errRes := handler(conn, &Request{ID: msg.ID, Method: msg.Method, Params: msg.Params}); res != nil || errRes != nil {
				if _, err := conn.WriteMsg(Response{ID: msg.ID, Result: res, Error: errRes}); err != nil {
					log.Fatalln("Errored while sending JSON-RPC response:", err)
				}
			}
		}
	} else { // if notification
		if handler, ok := r.notificationRoutes[msg.Method]; ok {
			handler(conn, &Notification{Method: msg.Method, Params: msg.Params})
		}
	}

}
