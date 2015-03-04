package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
)

// Listener accepts connections from devices.
type Listener struct {
	debug    bool
	listener net.Listener
}

// Listen creates a listener with the given PEM encoded X.509 certificate and the private key on the local network address laddr.
// Debug mode logs all server activity.
func Listen(cert, priv []byte, laddr string, debug bool) (*Listener, error) {
	c, err := x509.ParseCertificate(cert)
	p, err := x509.ParsePKCS1PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the certificate or the private key with error: %v", err)
	}

	pool := x509.NewCertPool()
	pool.AddCert(c)

	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  p,
	}

	config := tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientCAs:    pool,
	}

	if laddr == "" {
		laddr = "0.0.0.0:443"
	}

	listener, err := tls.Listen("tcp", laddr, &config)
	if err != nil {
		return nil, err
	}

	if debug {
		log.Printf("New CCS connection established with XMPP parameters: %+v\n", c)
	}

	return &Listener{
		debug:    debug,
		listener: listener,
	}, nil
}

// Accept waits for incoming connections and forwards incoming messages to handleMsg in a new goroutine.
// This function never returns, unless there is an error while accepting new connections.
func (l *Listener) Accept(handleMsg func(string)) error {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			log.Printf("error while accepting a new connection from a client: %v", err)
			return err // todo: it might not be appropriate to break the loop on all errors (like client disconnect during handshake)
		}

		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}

		tlscon, ok := conn.(*tls.Conn)
		if ok {
			state := tlscon.ConnectionState()
			sub := state.PeerCertificates[0].Subject
			log.Println(sub)
		}

		log.Printf("server: conn: echo %q\n", string(buf[:n]))
		n, err = conn.Write(buf[:n])

		n, err = conn.Write(buf[:n])
		log.Printf("server: conn: wrote %d bytes", n)

		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Println("server: conn: closed")
}
