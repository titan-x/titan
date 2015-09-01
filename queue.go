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
	conns *cmap.CMap      // user ID -> conn ID (reverse of what is stored in neptulon server)
	route *jsonrpc.Router // route to send messages through
}

// NewQueue creates a new queue object.
func NewQueue() Queue {
	return Queue{
		conns: cmap.New(),
	}
}

// AddConn adds a new connection ID for a given user.
func (q *Queue) AddConn(userID, connID string) {
	q.conns.Set(userID, connID)
}

// SendRequest attempts to send a request message to a given user.
// If user is not connected, message is queued.
func (q *Queue) SendRequest(userID string, request *jsonrpc.Request) {
	if connID, ok := q.conns.Get(userID); ok {
		// todo: use a client instance in sender as it already implements simplified sending, array sending functions
		q.route.SendRequest(connID.(string), request)
	} else {
		// todo: queue the message to userID for later delivery
	}
}
