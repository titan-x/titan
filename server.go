package main

import (
	"crypto/tls"
	"crypto/x509"
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
func (s *Server) Accept() error {
	return s.listener.Accept(handleMsg, handleDisconn)
}

// Stop stops a server instance gracefully, waiting for remaining data to be written on open connections.
func (s *Server) Stop() error {
	return nil
}

func handleMsg(conn *tls.Conn, session *Session, msg []byte) {
	if session.id == "" {
		auth(conn.ConnectionState().PeerCertificates, session, msg)
	}
}

func auth(peerCerts []*x509.Certificate, session *Session, msg []byte) {
	// client certificate authorization: certificate is verified by the listener instance so we trust it
	if len(peerCerts) > 0 {
		session.id = peerCerts[0].Subject.CommonName
		log.Printf("Client connected with client certificate subject: %+v", peerCerts[0].Subject)
	}

	// username/password authentication
	// todo: json/func Unmarshal(data []byte, v interface{}) error
}

func handleDisconn(conn *tls.Conn, session *Session) {

}
