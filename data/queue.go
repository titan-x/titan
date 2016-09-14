package data

// var queueLength = expvar.NewInt("queue-length") // todo: conn ID -> *Conn is tracked by neptulon already so we don't need to dupe this code
// var conns = expvar.NewInt("conns")

// Queue is a message queue for queueing and sending messages to users.
type Queue interface {
	// todo: buffered channels or basic locks or a concurrent multimap?
	// todo: at-least-once delivery relaxes things a bit for queueProcessor
	//
	// actually queue should not be interacted with directly, just like DB, it should be an interface
	// and server.send(userID) should use it automatically behind the scenes
}
