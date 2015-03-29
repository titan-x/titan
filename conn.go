package main

import (
	"crypto/tls"
	"time"
)

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

// ReadMsg waits for and reads the next message of the TLS connection, up to the given max message size. If max size is not defined,
// header limit of 4294967295 bytes (4GB) is used.
func (c *Conn) ReadMsg(maxSize int) (msg []byte, err error) {
	if c.err != nil {
		return nil, c.err
	}

	if err = c.conn.SetReadDeadline(time.Now().Add(time.Minute * 5)); err != nil {
		return
	}

	// first 4 bytes (uint32) is message length header with a maximum of 4294967295 (4GB) with a hard-cap defined by the parameter

	return
}
