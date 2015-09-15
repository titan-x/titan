package devastator

import (
	"log"

	"github.com/nbusy/cmap"
	"github.com/nbusy/neptulon/jsonrpc"
)

// Queue is a message queue for queueing and sending messages to users.
type Queue struct {
	conns *cmap.CMap                   // user ID -> conn ID
	route *jsonrpc.Router              // route to send messages through
	reqs  map[string]([]queuedRequest) // user ID -> []queuedRequest
}

// NewQueue creates a new queue object and registers the JSON-RPC router to be used for sending queued messages.
// Also attaches proper Neptulon middleware to initiate queue processing.
func NewQueue(server *jsonrpc.Server, r *jsonrpc.Router) Queue {
	q := Queue{
		conns: cmap.New(),
		route: r,
		reqs:  make(map[string]([]queuedRequest)),
	}

	server.ReqMiddleware(func(ctx *jsonrpc.ReqCtx) {
		q.SetConn(ctx.Conn.Data.Get("userid").(string), ctx.Conn.ID)
	})
	server.ResMiddleware(func(ctx *jsonrpc.ResCtx) {
		q.SetConn(ctx.Conn.Data.Get("userid").(string), ctx.Conn.ID)
	})
	server.NotMiddleware(func(ctx *jsonrpc.NotCtx) {
		q.SetConn(ctx.Conn.Data.Get("userid").(string), ctx.Conn.ID)
	})

	return q
}

// SetConn associates a user with a connection by ID.
// If there are pending messages for the user, they are started to be send immediately.
func (q *Queue) SetConn(userID, connID string) {
	if _, ok := q.conns.GetOk(userID); !ok {
		q.conns.Set(userID, connID)
		go q.processQueue(userID)
	}
}

// RemoveConn removes a user's associated connection ID.
func (q *Queue) RemoveConn(userID string) {
	q.conns.Delete(userID)
}

// AddRequest queues a request message to be sent to the given user.
func (q *Queue) AddRequest(userID string, method string, params interface{}, resHandler func(ctx *jsonrpc.ResCtx)) error {
	r := queuedRequest{Method: method, Params: params, ResHandler: resHandler}
	if rs, ok := q.reqs[userID]; ok {
		q.reqs[userID] = append(rs, r)
	} else {
		q.reqs[userID] = []queuedRequest{{Method: method, Params: params, ResHandler: resHandler}}
	}

	go q.processQueue(userID)
	return nil
}

type queuedRequest struct {
	Method     string
	Params     interface{}
	ResHandler func(ctx *jsonrpc.ResCtx)
}

// todo: prevent concurrent runs of processQueue or make []queuedRequest thread-safe
func (q *Queue) processQueue(userID string) {
	connID, ok := q.conns.GetOk(userID)
	if !ok {
		return
	}

	if reqs, ok := q.reqs[userID]; ok {
		for i, req := range reqs {
			if err := q.route.SendRequest(connID.(string), req.Method, req.Params, req.ResHandler); err != nil {
				log.Fatal(err) // todo: log fatal only in debug mode
			} else {
				reqs, reqs[len(reqs)-1] = append(reqs[:i], reqs[i+1:]...), queuedRequest{} // todo: this might not be needed if function is not a pointer val
			}
		}
	}
}
