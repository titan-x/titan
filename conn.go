package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const headerSize = 4

// Conn is a mobile client connection.
type Conn struct {
	conn         *tls.Conn
	isClient     bool
	maxMsgSize   int
	readDeadline time.Duration
}

// NewConn creates a new server-side connection object. Default values for maxMsgSize and readDeadline are
// 4294967295 bytes (4GB) and 300 seconds, respectively.
func NewConn(conn *tls.Conn, maxMsgSize int, readDeadline int) (*Conn, error) {
	if maxMsgSize == 0 {
		maxMsgSize = 4294967295
	}

	if readDeadline == 0 {
		readDeadline = 300
	}

	return &Conn{
		conn:         conn,
		maxMsgSize:   maxMsgSize,
		readDeadline: time.Second * time.Duration(readDeadline),
	}, nil
}

// Dial creates a new client side connection to a given network address with optional root CA and/or a client certificate (PEM encoded X.509 cert/key).
func Dial(addr string, rootCA []byte, clientCert []byte, clientCertKey []byte) (*Conn, error) {
	var roots *x509.CertPool
	var certs []tls.Certificate
	if rootCA != nil {
		roots = x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(rootCA)
		if !ok {
			return nil, errors.New("failed to parse root certificate")
		}
	}
	if clientCert != nil {
		tlsCert, err := tls.X509KeyPair(clientCert, clientCertKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the client certificate: %v", err)
		}
		certs = []tls.Certificate{tlsCert}
	}

	c, err := tls.Dial("tcp", addr, &tls.Config{RootCAs: roots, Certificates: certs})
	if err != nil {
		return nil, err
	}

	return NewConn(c, 0, 0)
}

// Write given message to the connection.
func (c *Conn) Write(msg *interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Failed to serialize the given message: %v", err)
	}

	n, err := c.conn.Write(data)
	if n != len(data) {
		return errors.New("Given message data length and sent bytes length did not match")
	}

	return err
}

// Read waits for and reads the next message of the TLS connection.
func (c *Conn) Read() (msg []byte, err error) {
	if err = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline)); err != nil {
		return
	}

	// first 4 bytes (uint32) is message length header with a maximum of 4294967295 bytes of message body (4GB) or the hard-cap defined by the user
	h := make([]byte, headerSize)
	n, err := c.conn.Read(h)
	if err != nil {
		return
	}
	if n != headerSize {
		return nil, fmt.Errorf("failed to read %v bytes message header, instead only read %v bytes", headerSize, n)
	}

	n = int(binary.LittleEndian.Uint32(h))
	r := 0
	msg = make([]byte, n)
	for r != n {
		for r != n {
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
	// todo: if session.err is nil, send a close req and wait ack then close? (or even wait for everything else to finish?)
	return c.conn.Close()
}
