package main

// Context is an incoming message context.
type Context struct {
	UserID         uint32 // should be in state
	IsDisconnected bool
	Error          string
	Conn           *Conn
	State          interface{}
}
