package inmem

import (
	"sync"

	"github.com/neptulon/neptulon"
	"github.com/titan-x/titan/data"
)

// Queue is a message queue for queueing and sending messages to users.
type Queue struct {
	senderFunc SenderFunc                // sender function to send and receive messages through
	conns      map[string]string         // user ID -> conn ID
	reqChans   map[string]queueProcessor // user ID -> queueProcessor
	mutex      sync.RWMutex
}

// NewQueue creates a new queue object.
func NewQueue(senderFunc SenderFunc) *Queue {
	q := Queue{
		senderFunc: senderFunc,
		conns:      make(map[string]string),
		reqChans:   make(map[string]queueProcessor),
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

type queueProcessor struct {
	reqChan  chan queuedRequest
	quitChan chan bool
}

func (q *Queue) getQueueProcessor(userID string) queueProcessor {
	c, ok := q.reqChans[userID]
	if !ok {
		c = queueProcessor{reqChan: make(chan queuedRequest, 5000), quitChan: make(chan bool)}
		q.reqChans[userID] = c
	}
	return c
}

// Middleware registers a queue middleware to register user/connection IDs
// for connecting users (upon their first incoming-message).
func (q *Queue) Middleware(ctx *neptulon.ReqCtx) error {
	q.mutex.Lock()
	userID := ctx.Conn.Session.Get("userid").(string)
	connID := ctx.Conn.ID

	// start queue gorutine once per connection
	if _, ok := q.conns[userID]; !ok {
		q.conns[userID] = connID
		data.UserCount.Add(1)
		go q.processQueue(connID, q.getQueueProcessor(userID))
	}
	q.mutex.Unlock()

	return ctx.Next()
}

// RemoveConn removes a user's associated connection ID.
func (q *Queue) RemoveConn(userID string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	c := q.getQueueProcessor(userID)
	c.quitChan <- true
	delete(q.conns, userID)
	data.UserCount.Add(-1)
}

// AddRequest queues a request message to be sent to the given user.
func (q *Queue) AddRequest(userID string, method string, params interface{}, resHandler func(ctx *neptulon.ResCtx) error) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	data.QueueLength.Add(1)
	r := queuedRequest{Method: method, Params: params, ResHandler: resHandler}
	c := q.getQueueProcessor(userID)
	c.reqChan <- r

	return nil
}

func (q *Queue) processQueue(connID string, processor queueProcessor) {
	for {
		select {
		case req := <-processor.reqChan:
			if _, err := q.senderFunc(connID, req.Method, req.Params, req.ResHandler); err != nil {
				// write the request back to buffered channel and continue until quit
				processor.reqChan <- req
				continue
			}
			data.QueueLength.Add(-1)
		case <-processor.quitChan:
			return
		}
	}
}
