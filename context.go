package main

// Context is a connection context.
type Context struct {
	UserID       uint32 // should be in state? as user management is not a listener package duty..
	Disconnected bool
	err          string
	Conn         *Conn
	State        interface{}
}

// SetError sets a connection error that will be written to the connection before it is closed right before the next read cycle.
func (c *Context) SetError(err string) {
	// todo use mutex
	c.err = err
}
