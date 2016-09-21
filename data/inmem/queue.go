package inmem

import (
	"github.com/neptulon/neptulon"
	"github.com/titan-x/titan/data"
)

// Queue is a message queue for queueing and sending messages to users.
type Queue struct {
	senderFunc SenderFunc           // sender function to send and receive messages through
	conns      map[string]string    // user ID -> conn ID
	reqChans   map[string]queueChan // user ID -> queueProcessor

	// worker communication channels
	middlewareChan chan middlewareChan
	remUserChan    chan string
	addReqChan     chan addReqChan
	delQueueChan   chan string
}

// NewQueue creates a new queue object.
func NewQueue(senderFunc SenderFunc) *Queue {
	q := Queue{
		senderFunc: senderFunc,
		conns:      make(map[string]string),
		reqChans:   make(map[string]queueChan),

		middlewareChan: make(chan middlewareChan, 5000),
		remUserChan:    make(chan string, 5000),
		addReqChan:     make(chan addReqChan, 5000),
		delQueueChan:   make(chan string, 5000),
	}

	go q.worker()

	return &q
}

// SenderFunc is a function for sending messages over connections ID.
type SenderFunc func(connID string, method string, params interface{}, resHandler func(ctx *neptulon.ResCtx) error) (reqID string, err error)

type queuedReq struct {
	Method     string
	Params     interface{}
	ResHandler func(ctx *neptulon.ResCtx) error
}

type queueChan struct {
	req  chan queuedReq
	quit chan bool
}

func (q *Queue) getQueueChan(userID string) queueChan {
	c, ok := q.reqChans[userID]
	if !ok {
		c = queueChan{req: make(chan queuedReq, 5000), quit: make(chan bool, 1)}
		q.reqChans[userID] = c
	}
	return c
}

// Middleware registers a queue middleware to register user/connection IDs
// for connecting users (upon their first incoming-message).
func (q *Queue) Middleware(ctx *neptulon.ReqCtx) error {
	q.middlewareChan <- middlewareChan{userID: ctx.Conn.Session.Get("userid").(string), connID: ctx.Conn.ID}
	return ctx.Next()
}

// RemoveConn removes a user's associated connection ID.
func (q *Queue) RemoveConn(userID string) {
	q.remUserChan <- userID
}

// AddRequest queues a request message to be sent to the given user.
func (q *Queue) AddRequest(userID string, method string, params interface{}, resHandler func(ctx *neptulon.ResCtx) error) error {
	q.addReqChan <- addReqChan{userID: userID, queuedReq: queuedReq{Method: method, Params: params, ResHandler: resHandler}}
	return nil
}

func (q *Queue) processQueue(qc queueChan, userID, connID string) {
	errc := 0 // protect against infinite retry loop

	for {
		select {
		case req := <-qc.req:
			_, err := q.senderFunc(connID, req.Method, req.Params, req.ResHandler)

			if err != nil {
				errc++
				qc.req <- req
				if errc > 10 {
					return
				}
				continue
			}

			data.QueueLength.Add(-1)
			errc = 0

		case <-qc.quit:
			if len(qc.req) == 0 {
				q.delQueueChan <- userID
			}
			return
		}
	}
}
