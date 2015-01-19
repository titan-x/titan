package main

import "crypto/tls"

// Conn is a mobile client connection.
type Conn struct {
	conn tls.Conn

	// user -> id (user or chat id) -> message
	// delivery status -> user
	// read status -> user
}
