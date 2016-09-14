package data

import (
	"expvar"

	"github.com/neptulon/neptulon"
)

var queueLength = expvar.NewInt("queue-length")
var usersCount = expvar.NewInt("users") // authenticated users only

// Queue is a message queue for queueing and sending messages to users.
type Queue interface {
	// todo: buffered channels or basic locks or a concurrent multimap?
	// todo: at-least-once delivery relaxes things a bit for queueProcessor
	//
	// actually queue should not be interacted with directly, just like DB, it should be an interface
	// and server.send(userID) should use it automatically behind the scenes

	Middleware(ctx *neptulon.ReqCtx) error
	SetConn(userID, connID string)
	RemoveConn(userID string)
	AddRequest(userID string, method string, params interface{}, resHandler func(ctx *neptulon.ResCtx) error) error
}
