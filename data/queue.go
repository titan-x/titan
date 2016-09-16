package data

import (
	"expvar"

	"github.com/neptulon/neptulon"
)

// Queue is a message queue for queueing and sending messages to users.
type Queue interface {
	Middleware(ctx *neptulon.ReqCtx) error
	SetConn(userID, connID string)
	RemoveConn(userID string)
	AddRequest(userID string, method string, params interface{}, resHandler func(ctx *neptulon.ResCtx) error) error
}

// QueueLength is the total request queue for all users combined.
// This should be handled by the implementing struct.
var QueueLength = expvar.NewInt("queue-length")

// UserCount is the total authenticated live user count.
// This should be handled by the implementing struct.
var UserCount = expvar.NewInt("users")
