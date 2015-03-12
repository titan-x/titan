package main

import (
	"crypto/tls"
	"log"
)

// Server wraps a listener instance and registers default connection and message handlers with the listener.
type Server struct {
	debug    bool
	listener *Listener
}

// NewServer creates and returns a new server instance with a listener created using given parameters.
func NewServer(cert, privKey []byte, laddr string, debug bool) (*Server, error) {
	l, err := Listen(cert, privKey, laddr, debug)
	if err != nil {
		return nil, err
	}

	return &Server{
		debug:    debug,
		listener: l,
	}, nil
}

// Accept accepts connections on the internal listener and handles connections with registered onnection and message handlers.
// This function blocks and never returns, unless there is an error while accepting a new connection.
func (s *Server) Accept() {
	s.listener.Accept(func(conn *tls.Conn, session *Session, msg []byte) {
		// wg.Add(1)
		// defer wg.Done()
		log.Printf("Incoming message to listener from a client: %v", string(msg))

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) > 0 {
			log.Printf("Client connected with client certificate subject: %v", certs[0].Subject)
		}
	}, func(conn *tls.Conn, session *Session) {
	})
}

// Stop stops a server instance gracefully, waiting for remaining data to be written on open connections.
func (s *Server) Stop() error {
	return nil
}
