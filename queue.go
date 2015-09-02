package devastator

import (
	"github.com/nbusy/cmap"
	"github.com/nbusy/neptulon/jsonrpc"
)

// Queue is a message queue for queueing and sending messages to users.
type Queue struct {
	conns *cmap.CMap                   // user ID -> conn ID
	route *jsonrpc.Router              // route to send messages through
	reqs  map[string]([]queuedRequest) // user ID -> []queuedRequest
}

// NewQueue creates a new queue object.
func NewQueue(r *jsonrpc.Router) Queue {
	return Queue{
		conns: cmap.New(),
		route: r,
		reqs:  make(map[string]([]queuedRequest)),
	}
}

// SetConn associates a user with a connection by ID.
// If there are pending messages for the user, they are started to be send immediately.
func (q *Queue) SetConn(userID, connID string) {
	q.conns.Set(userID, connID)
	// todo: trigger processing
}

// RemoveConn removes a user's associated connection ID.
func (q *Queue) RemoveConn(userID string) {
	q.conns.Delete(userID)
}

// AddRequest queues a request message to be sent to the given user.
func (q *Queue) AddRequest(userID string, method string, params interface{}, resHandler func(ctx *jsonrpc.ResCtx)) error {
	if connID, ok := q.conns.Get(userID); ok {
		if err := q.route.SendRequest(connID.(string), method, params, resHandler); err != nil {
			return err
		}
	} else {
		r := queuedRequest{Method: method, Params: params}
		if rs, ok := q.reqs[userID]; ok {
			q.reqs[userID] = append(rs, r)
		} else {
			q.reqs[userID] = []queuedRequest{{Method: method, Params: params}}
		}
	}

	return nil
}

type queuedRequest struct {
	Method string
	Params interface{}
}
