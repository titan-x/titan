package main

import "crypto/tls"

// Conn is a mobile client connection.
type Conn struct {
	UserID uint32
	Error  string
	Data   interface{}
	conn   *tls.Conn
}

// SendMessage sends a message to the connected mobile client.
func (c *Conn) SendMessage() error {
	return nil
}

// SendNotification sends a notification to the connected mobile client.
func (c *Conn) SendNotification() error {
	return nil
}
