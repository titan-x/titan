package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

const headerSize = 4

// Conn is a mobile client connection.
type Conn struct {
	conn              *tls.Conn
	maxMsgSize        int
	readWriteDeadline time.Duration
}

// NewConn creates a new server-side connection object. Default values for maxMsgSize and readWriteDeadline are
// 4294967295 bytes (4GB) and 300 seconds, respectively.
func NewConn(conn *tls.Conn, maxMsgSize int, readWriteDeadline int) *Conn {
	if maxMsgSize == 0 {
		maxMsgSize = 4294967295
	}

	if readWriteDeadline == 0 {
		readWriteDeadline = 300
	}

	return &Conn{
		conn:              conn,
		maxMsgSize:        maxMsgSize,
		readWriteDeadline: time.Second * time.Duration(readWriteDeadline),
	}
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

	return NewConn(c, 0, 0), nil
}

// Write given message to the connection.
func (c *Conn) Write(msg *interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Failed to serialize the given message: %v", err)
	}

	if err = c.conn.SetReadDeadline(time.Now().Add(c.readWriteDeadline)); err != nil {
		return err
	}

	n, err := c.conn.Write(data)
	if n != len(data) {
		return errors.New("Given message data length and sent bytes length did not match")
	}

	return err
}

// Read waits for and reads the next message of the TLS connection.
func (c *Conn) Read() (msg []byte, err error) {
	if err = c.conn.SetReadDeadline(time.Now().Add(c.readWriteDeadline)); err != nil {
		return
	}

	// read the content length header
	reader := bufio.NewReader(c.conn)
	line, err := reader.ReadSlice('\n')
	if err != nil {
		if err == io.EOF {
			return
		}

		log.Fatalln("Client read error:", err)
	}

	// calculate the content length
	n, err := strconv.Atoi(string(line[:len(line)-1]))
	if err != nil || n == 0 {
		log.Fatalln("Client read error: invalid content lenght header sent or content lenght mismatch:", err)
	}

	// read the message content
	msg = make([]byte, n)
	total := 0
	for total != n {
		// todo: log here in case it gets stuck, or there is a dos attack, pumping up cpu usage!
		i, err := reader.Read(msg[total:])
		if err != nil {
			log.Fatalln("Error while reading incoming message:", err)
			break
		}
		total += i
	}
	if err != nil {
		log.Fatalln("Error while reading incoming message:", err)
	}

	return
}

// Close closes a connection.
func (c *Conn) Close() error {
	// todo: if session.err is nil, send a close req and wait ack then close? (or even wait for everything else to finish?)
	return c.conn.Close()
}
