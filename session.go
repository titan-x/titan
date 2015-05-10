package main

import "sync"

// Session is a session data store for connections.
type Session struct {
	UserID       uint32
	Error        error
	Disconnected bool
	data         map[string]interface{}
	mutex        sync.RWMutex
}

// Set stores a value for a given key in the session. This method is thread safe.
func (s *Session) Set(key string, val interface{}) {
	s.mutex.Lock()
	s.data[key] = val
	s.mutex.Unlock()
}

// Get retrieves a value for a given key in the session. This method is thread safe.
func (s *Session) Get(key string) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key]
}
