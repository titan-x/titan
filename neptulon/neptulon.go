// Package neptulon is a socket framework with middleware support.
package neptulon

import (
	"net/http"
	"sync"
)

// Neptulon framework entry point.
type Neptulon struct {
	debug       bool
	err         error
	listener    *Listener
	mutex       sync.Mutex
	middlewares []*func(ctx Context) (response interface{})
}

// Handle registers a new middleware to handle incoming messages.
func (n *Neptulon) Handle() {}

type handler func(w http.ResponseWriter, r *http.Request) error

func handle(handlers ...handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			err := handler(w, r)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}
	})
}
