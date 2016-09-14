package titan

import "expvar"

var queueLength = expvar.NewInt("queue-length")
var conns = expvar.NewInt("conns")

// Queue2 is a message queue for queueing and sending messages to users.
type Queue2 struct {
	// todo: buffered channels or basic locks or a concurrent multimap?
	// todo: at-least-once delivery relaxes things a bit for queueProcessor
	//
	// actually queue should not be interacted with directly, just like DB, it should be an interface
	// and server.send(userID) should use it automatically behind the scenes
}
