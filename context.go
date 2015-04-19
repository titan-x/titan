package main

// Context is a connection context.
type Context struct {
	UserID         uint32 // should be in state
	IsDisconnected bool
	error          string
	Conn           *Conn
	State          interface{}
}

// SetError sets a connection error that will be written to the connection before it is closed right before the next read cycle.
func (c *Context) SetError(err string) {
	c.error = err
}
