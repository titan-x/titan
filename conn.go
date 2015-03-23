package main

import "crypto/tls"

// Conn is a mobile client connection.
type Conn struct {
	conn    *tls.Conn
	session interface{}

	// user -> id (user or chat id) -> message
	// delivery status -> user
	// read status -> user
}

// SendMessage sends a message to the connected mobile client.
func (c *Conn) SendMessage() error {
	return nil
}

// SendNotification sends a notification to the connected mobile client.
func (c *Conn) SendNotification() error {
	return nil
}
