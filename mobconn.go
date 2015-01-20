package main

import "crypto/tls"

// MobConn is a mobile client connection.
type MobConn struct {
	conn tls.Conn

	// user -> id (user or chat id) -> message
	// delivery status -> user
	// read status -> user
}

// SendMessage sends a message to a connected mobile client.
func (c *MobConn) SendMessage() error {
	return nil
}

// SendNotification sends a notification to a connected mobile client.
func (c *MobConn) SendNotification() error {
	return nil
}
