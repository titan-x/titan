package main

import "sync"

// Session is a session data store for connections. All operations on session are thread safe.
type Session struct {
	Error        error // todo: use mutex or remove these in favor of data. could stay if these end up being listener internal only
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
