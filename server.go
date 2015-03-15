package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

var conns = make(map[uint32]*tls.Conn)

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
	if session.UserID == 0 {
		userID, err := auth(conn.ConnectionState().PeerCertificates, msg)
		if err != nil {
			session.Error = fmt.Sprintf("Cannot parse client message or method mismatched: %v", err)
		}
		session.UserID = userID
		conns[userID] = conn
		// todo: ack auth message, start sending other queued messages one by one
		// can have 2 approaches here
		// 1. users[id].send(...) & users[id].queue(...)
		// 2. conn.write(...) && queue[id].conn = ...
		return
	}

	// process the message and queue a reply if necessary
}

// auth handles Google+ sign-in and client certificate authentication.
func auth(peerCerts []*x509.Certificate, msg []byte) (userID uint32, err error) {
	// client certificate authorization: certificate is verified by the TLS listener instance so we trust it
	if len(peerCerts) > 0 {
		idstr := peerCerts[0].Subject.CommonName
		uid64, err := strconv.ParseUint(idstr, 10, 32)
		if err != nil {
			return 0, err
		}
		userID = uint32(uid64)
		log.Printf("Client connected with client certificate subject: %+v", peerCerts[0].Subject)
		return userID, nil
	}

	// Google+ authentication
	var req ReqMsg
	if err = json.Unmarshal(msg, &req); err != nil {
		return
	}

	switch req.Method {
	case "auth.token":
		var token string
		if err = json.Unmarshal(req.Params, &token); err != nil {
			return
		}
		// assume that token = user ID for testing
		uid64, err := strconv.ParseUint(token, 10, 32)
		if err != nil {
			return 0, err
		}
		userID = uint32(uid64)
		return userID, nil
	case "auth.google":
		// todo: ping google, get user info, save user info in DB, generate and return permanent jwt token (or should this part be NBusy's business?)
		return
	default:
		return 0, errors.New("initial unauthenticated request should be in the 'auth.xxx' form")
	}
}

func handleDisconn(conn *tls.Conn, session *Session) {
	delete(conns, session.UserID)
}
