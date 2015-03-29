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
		return nil, c.err
	}

	// first 4 bytes (uint32) is message length header with a maximum of 4294967295 (4GB) with a hard-cap defined by configuration

	return
}
