package neptulon

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// Conn is a mobile client connection.
type Conn struct {
	Session           Session
	conn              *tls.Conn
	headerSize        int
	maxMsgSize        int
	readWriteDeadline time.Duration
	debug             bool
}

// NewConn creates a new server-side connection object.
// Default values for headerSize, maxMsgSize, and readWriteDeadline are 4 bytes, 4294967295 bytes (4GB), and 300 seconds, respectively.
// Debug mode logs all raw TCP communication.
func NewConn(conn *tls.Conn, headerSize, maxMsgSize, readWriteDeadline int, debug bool) (*Conn, error) {
	if headerSize == 0 {
		headerSize = 4
	}
	if maxMsgSize == 0 {
		maxMsgSize = 4294967295
	}
	if readWriteDeadline == 0 {
		readWriteDeadline = 300
	}

	id, err := getID()
	if err != nil {
		return nil, err
	}

	return &Conn{
		Session:           Session{ID: id},
		conn:              conn,
		headerSize:        headerSize,
		maxMsgSize:        maxMsgSize,
		readWriteDeadline: time.Second * time.Duration(readWriteDeadline),
		debug:             debug,
	}, nil
}

// Dial creates a new client side connection to a given network address with optional root CA and/or a client certificate (PEM encoded X.509 cert/key).
// Debug mode logs all raw TCP communication.
func Dial(addr string, rootCA []byte, clientCert []byte, clientCertKey []byte, debug bool) (*Conn, error) {
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

	// todo: dial timeout like that of net.Conn.DialTimeout
	c, err := tls.Dial("tcp", addr, &tls.Config{RootCAs: roots, Certificates: certs})
	if err != nil {
		return nil, err
	}

	return NewConn(c, 0, 0, 0, debug)
}

// Read waits for and reads the next incoming message from the TLS connection.
func (c *Conn) Read() (n int, msg []byte, err error) {
	if err = c.conn.SetReadDeadline(time.Now().Add(c.readWriteDeadline)); err != nil {
		return
	}

	// read the content length header
	h := make([]byte, c.headerSize)
	n, err = c.conn.Read(h)
	if err != nil {
		return
	}
	if n != c.headerSize {
		err = fmt.Errorf("expected to read header size %v bytes but instead read %v bytes", c.headerSize, n)
		return
	}

	// calculate the content length
	n = readHeaderBytes(h)

	// read the message content
	msg = make([]byte, n)
	total := 0
	for total < n {
		// todo: log here in case it gets stuck, or there is a dos attack, pumping up cpu usage!
		i, err := c.conn.Read(msg[total:])
		if err != nil {
			err = fmt.Errorf("errored while reading incoming message: %v", err)
			break
		}
		total += i
	}
	if total != n {
		err = fmt.Errorf("expected to read %v bytes instead read %v bytes", n, total)
	}

	if c.debug {
		log.Println("Incoming message:", string(msg))
	}

	return
}

// Write writes given message to the connection.
func (c *Conn) Write(msg []byte) (n int, err error) {
	l := len(msg)
	h := makeHeaderBytes(l, c.headerSize)

	// write the header
	n, err = c.conn.Write(h)
	if err != nil {
		return
	}
	if n != c.headerSize {
		err = fmt.Errorf("expected to write %v bytes but only wrote %v bytes", l, n)
	}

	// write the body
	// todo: do we need a loop? bufio uses a loop but it might be due to buff length limitation
	n, err = c.conn.Write(msg)
	if err != nil {
		return
	}
	if n != l {
		err = fmt.Errorf("expected to write %v bytes but only wrote %v bytes", l, n)
	}

	return
}

// Close closes a connection.
func (c *Conn) Close() error {
	// todo: if session.err is nil, send a close req and wait ack then close? (or even wait for everything else to finish?)
	return c.conn.Close()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// ConnectionState returns basic TLS details about the connection.
func (c *Conn) ConnectionState() tls.ConnectionState {
	return c.conn.ConnectionState()
}

func makeHeaderBytes(h, size int) []byte {
	b := make([]byte, size)
	binary.LittleEndian.PutUint32(b, uint32(h))
	return b
}

func readHeaderBytes(h []byte) int {
	return int(binary.LittleEndian.Uint32(h))
}
