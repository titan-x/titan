package main

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"time"
)

const headerSize = 4

// Conn is a mobile client connection.
type Conn struct {
	UserID       uint32
	conn         *tls.Conn
	err          error
	maxMsgSize   int
	readDeadline time.Duration
	header       []byte
}

// NewConn creates and returns a new connection object. Default values for maxMsgSize and readDeadline are 4294967295 bytes (4GB)
// and 300 seconds, respectively.
func NewConn(conn *tls.Conn, maxMsgSize int, readDeadline int) (*Conn, error) {
	if maxMsgSize == 0 {
		maxMsgSize = 4294967295
	}

	if readDeadline == 0 {
		readDeadline = 300
	}

	return &Conn{
		header:       make([]byte, headerSize), // todo: use a regular byte array rather than a slice?
		conn:         conn,
		maxMsgSize:   maxMsgSize,
		readDeadline: time.Second * time.Duration(readDeadline),
	}, nil
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

	if err = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline)); err != nil {
		return
	}

	// first 4 bytes (uint32) is message length header with a maximum of 4294967295 bytes of message body (4GB) or the hard-cap defined by the user
	n, err := c.conn.Read(c.header)
	if err != nil {
		return
	}
	if n != headerSize {
		return nil, fmt.Errorf("failed to read %v bytes message header, instead only read %v bytes", headerSize, n)
	}

	t := int(binary.LittleEndian.Uint32(c.header))
	r := 0
	msg = make([]byte, t)
	for {
		for r != t {
			i, err := c.conn.Read(msg[r:])
			if err != nil {
				return nil, fmt.Errorf("errored while reading incoming message: %v", err)
			}
			r += i
		}
	}

	return
}

// Close closes a connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}
