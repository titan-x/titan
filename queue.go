package main

// Queue is a message queue for connected devices.
// Messages are mapped from user ID to []interface{} array which may contain request, response, or notification messages.
type Queue map[uint32][]interface{}

// Send
func (q *Queue) Send() {}
