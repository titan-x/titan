// Package neptulon is a socket framework with middleware support.
package neptulon

import "net/http"

var middlewares []*func(ctx Context) (response interface{})

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
