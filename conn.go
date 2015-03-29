package main

import "crypto/tls"

// Conn is a mobile client connection.
type Conn struct {
	UserID uint32
	conn   *tls.Conn
	err    error
}

// SendMsg sends a message to the connected mobile client.
func (c *Conn) SendMsg(msg *interface{}) error {
	return nil
}

// ReadMsg waits for and reads the next message of the TLS connection.
func (c *Conn) ReadMsg() (msg []byte, err error) {
	if c.err != nil {
		// todo: send error message to user, log the error, and close the conn and return
		return nil, c.err
	}

	return
}
