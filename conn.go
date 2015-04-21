package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

// Conn is a mobile client connection.
type Conn struct {
	conn              *tls.Conn
	headerSize        int
	maxMsgSize        int
	readWriteDeadline time.Duration
}

// NewConn creates a new server-side connection object.
// Default values for headerSize, maxMsgSize, and readWriteDeadline
// are 4 bytes, 4294967295 bytes (4GB), and 300 seconds, respectively.
func NewConn(conn *tls.Conn, headerSize, maxMsgSize, readWriteDeadline int) *Conn {
	if headerSize == 0 {
		headerSize = 4
	}
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

	return NewConn(c, 0, 0, 0), nil
}

// ReadMsg waits for and reads the next incoming message from the TLS connection and deserializes it into the given message object.
func (c *Conn) ReadMsg(msg interface{}) (n int, err error) {
	n, data, err := c.Read()
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, msg); err != nil {
		return
	}

	return
}

// Read waits for and reads the next incoming message from the TLS connection.
func (c *Conn) Read() (n int, msg []byte, err error) {
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
	if n, err = strconv.Atoi(string(line[:len(line)-1])); err != nil || n == 0 {
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

// WriteMsg serializes and writes given message to the connection with appropriate header.
func (c *Conn) WriteMsg(msg *interface{}) (n int, err error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize the given message: %v", err)
	}

	return c.Write(data)
}

// Write writes given message to the connection.
func (c *Conn) Write(msg []byte) (n int, err error) {
	l := strconv.Itoa(len(msg))
	h := append([]byte(l), []byte("\n")...)
	msg = append(h, msg...)

	return c.conn.Write(msg)
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// ConnectionState returns basic TLS details about the connection.
func (c *Conn) ConnectionState() tls.ConnectionState {
	return c.conn.ConnectionState()
}

// Close closes a connection.
func (c *Conn) Close() error {
	// todo: if session.err is nil, send a close req and wait ack then close? (or even wait for everything else to finish?)
	return c.conn.Close()
}

func makeHeaderBytes(i int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

func readHeaderBytes(h []byte) int {
	return int(binary.LittleEndian.Uint32(h))
}
