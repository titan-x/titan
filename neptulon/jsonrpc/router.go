package jsonrpc

import "github.com/nbusy/devastator/neptulon"

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

func (r *Router) middleware(conn *neptulon.Conn, msg *Message) (result interface{}, err *ResError) {
	// if not request or notification don't handle it
	if msg.Method == "" {
		return nil, nil
	}

	// if request
	if msg.ID != "" {
		if handler, ok := r.requestRoutes[msg.Method]; ok {
			if res, resErr := handler(conn, &Request{ID: msg.ID, Method: msg.Method, Params: msg.Params}); res != nil || resErr != nil {
				return res, resErr
			}
		}
	} else { // if notification
		if handler, ok := r.notificationRoutes[msg.Method]; ok {
			handler(conn, &Notification{Method: msg.Method, Params: msg.Params})
			// todo: need to return something to prevent deeper handlers to further handle this request (i.e. not found handler logging not found warning)
		}
	}

	return nil, nil
}
