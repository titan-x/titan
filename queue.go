package devastator

import (
	"github.com/nbusy/cmap"
	"github.com/nbusy/neptulon/jsonrpc"
)

// todo: type queue { conns cmap.CMap; route Route; queue Queue }
//         .AddConn() { // trigger queue send inside cert_auth as long as conn is available }
//         .Disconn() { }

// todo2: this now looks like the old Sender so it might be logical to rename this to Sender and implement separate Queue

// Queue is a message queue for queueing and sending messages to users.
type Queue struct {
	conns *cmap.CMap                      // user ID -> conn ID
	route *jsonrpc.Router                 // route to send messages through
	reqs  map[string]([]*jsonrpc.Request) // user ID -> []*jsonrpc.Request
}

// NewQueue creates a new queue object.
func NewQueue() Queue {
	return Queue{
		conns: cmap.New(),
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
func (q *Queue) AddRequest(userID string, request *jsonrpc.Request, resHandler func(ctx *jsonrpc.ResCtx)) {
	if connID, ok := q.conns.Get(userID); ok {
		q.route.SendRequest(connID.(string), request, resHandler)
	} else {
		// q.reqs[userID] = append(q.reqs[userID], request)
	}
}
