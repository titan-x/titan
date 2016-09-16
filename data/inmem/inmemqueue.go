package inmem

import (
	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon"
	"github.com/titan-x/titan/data"
)

// Queue is a message queue for queueing and sending messages to users.
type Queue struct {
	senderFunc SenderFunc // sender function to send and receive messages through
	conns      *cmap.CMap // user ID -> conn ID
	reqChans   *cmap.CMap // user ID -> make(chan *queuedRequest, 100)
}

// NewQueue creates a new queue object.
func NewQueue(senderFunc SenderFunc) *Queue {
	q := Queue{
		conns:    cmap.New(),
		reqChans: cmap.New(),
	}

	return &q
}

// SenderFunc is a function for sending messages over connections ID.
type SenderFunc func(connID string, method string, params interface{}, resHandler func(ctx *neptulon.ResCtx) error) (reqID string, err error)

type queuedRequest struct {
	Method     string
	Params     interface{}
	ResHandler func(ctx *neptulon.ResCtx) error
}

// Middleware registers a queue middleware to register user/connection IDs
// for connecting users (upon their first incoming-message).
func (q *Queue) Middleware(ctx *neptulon.ReqCtx) error {
	uid := ctx.Conn.Session.Get("userid")
	if _, ok := q.conns.GetOk(uid); !ok {
		q.SetConn(uid.(string), ctx.Conn.ID)
	}

	return ctx.Next()
}

// SetConn associates a user with a connection by ID.
// If there are pending messages for the user, they are started to be send immediately.
func (q *Queue) SetConn(userID, connID string) {
	if _, ok := q.conns.GetOk(userID); !ok {
		q.conns.Set(userID, connID)
		data.UserCount.Add(1)
		go q.processQueue()
	}
}

// RemoveConn removes a user's associated connection ID.
func (q *Queue) RemoveConn(userID string) {
	q.conns.Delete(userID)
	data.UserCount.Add(-1)
}

// AddRequest queues a request message to be sent to the given user.
func (q *Queue) AddRequest(userID string, method string, params interface{}, resHandler func(ctx *neptulon.ResCtx) error) error {
	data.QueueLength.Add(1)
	// r := queuedRequest{Method: method, Params: params, ResHandler: resHandler}

	return nil
}

func (q *Queue) processQueue() {
	// tood: data.QueueLength.Add(-1) per successful request
}
