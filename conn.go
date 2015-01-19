package main

import "crypto/tls"

// Conn is a mobile client connection.
type Conn struct {
	conn tls.Conn

	// user -> id (user or chat id) -> message
	// delivery status -> user
	// read status -> user
}

// SendMessage sends a message to a connected mobile client.
func (t *Conn) SendMessage() error {
	return nil
}

// SendNotification sends a notification to a connected mobile client.
func (t *Conn) SendNotification() error {
	return nil
}
