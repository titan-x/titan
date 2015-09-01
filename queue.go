package devastator

import (
	"github.com/nbusy/cmap"
	"github.com/nbusy/neptulon"
	"github.com/nbusy/neptulon/jsonrpc"
)

// todo: type queue { conns cmap.CMap; route Route; queue Queue }
//         .AddConn() { // trigger queue send inside cert_auth as long as conn is available }
//         .Disconn() { }

// todo2: this now looks like the old Sender so it might be logical to rename this to Sender and implement separate Queue

// Queue is a message queue for queueing and sending messages to users.
type Queue struct {
	conns *cmap.CMap      // user ID -> *neptulon.Conn
	route *jsonrpc.Router // route to send messages through
}

// NewQueue creates a new queue object.
func NewQueue() Queue {
	return Queue{
		conns: cmap.New(),
	}
}

// SetConn sets the active connection for the given user.
// If there are pending messages for the user, they start to be sent immediately.
func (q *Queue) SetConn(userID string, conn *neptulon.Conn) {
	q.conns.Set(userID, conn)
}

// Disconn releases dicsonnected user's connection object.
func (q *Queue) Disconn(userID string) {
	q.conns.Delete(userID)
}

// AddRequest queues a request message to be sent to the given user.
func (q *Queue) AddRequest(userID string, request *jsonrpc.Request) {
	if connID, ok := q.conns.Get(userID); ok {
		// todo: use a client instance in sender as it already implements simplified sending, array sending functions
		q.route.SendRequest(connID.(string), request)
	} else {
		// todo: queue the message to userID for later delivery
	}
}
