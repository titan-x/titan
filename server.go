package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
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

// handleMsg handles incoming client messages.
func handleMsg(conn *tls.Conn, session *Session, msg []byte) {
	// authenticate the session if not already done
	if session.UserID == "" {
		userID, err := auth(conn.ConnectionState().PeerCertificates, msg)
		if err != nil {
			session.Error = fmt.Sprintf("Cannot parse client message: %v", err)
		}
		session.UserID = userID
	}

	// todo: session is authenticated and we have user ID now so associate user ID with session in a go map (var users = make(map[uint32]User) maybe??)
}

// auth handles classical username/password and client certificate based authentication.
func auth(peerCerts []*x509.Certificate, msg []byte) (userID string, err error) {
	// client certificate authorization: certificate is verified by the listener instance so we trust it
	if len(peerCerts) > 0 {
		userID = peerCerts[0].Subject.CommonName
		log.Printf("Client connected with client certificate subject: %+v", peerCerts[0].Subject)
	}

	// username/password authentication
	var req ReqMsg
	err = json.Unmarshal(msg, req)
	if err != nil {
		return
	}

	return
}

func handleDisconn(conn *tls.Conn, session *Session) {

}
