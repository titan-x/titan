package main

import "sync"

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

// Session is a generic session data store for client handlers. All operations on session are thread safe.
type Session struct {
	UserID       uint32
	Error        error
	Disconnected bool
	data         map[string]interface{}
	mutex        sync.RWMutex
}

// Set stores a value for a given key in the session.
func (s *Session) Set(key string, val interface{}) {
	s.mutex.Lock()
	s.data[key] = val
	s.mutex.Unlock()
}

// Get retrieves a value for a given key in the session.
func (s *Session) Get(key string) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key]
}
