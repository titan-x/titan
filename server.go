package main

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
)

var users = make(map[uint32]*User)

// Server wraps a listener instance and registers default connection and message handlers with the listener.
type Server struct {
	debug    bool
	err      error
	listener *Listener
	acceptwg *sync.WaitGroup
	connwg   *sync.WaitGroup
	reqwg    *sync.WaitGroup
	mutex    sync.Mutex
}

// NewServer creates and returns a new server instance with a listener created using given parameters.
// Debug mode dumps raw TCP data to stderr (log.Println() default).
func NewServer(cert, privKey []byte, laddr string, debug bool) (*Server, error) {
	connwg := new(sync.WaitGroup)
	reqwg := new(sync.WaitGroup)

	l, err := Listen(cert, privKey, laddr, connwg, reqwg, debug)
	if err != nil {
		return nil, err
	}

	return &Server{
		debug:    debug,
		listener: l,
		connwg:   connwg,
		reqwg:    reqwg,
		mutex:    sync.Mutex{},
	}, nil
}

// Start starts accepting connections on the internal listener and handles connections with registered onnection and message handlers.
// This function blocks and never returns, unless there is an error while accepting a new connection.
func (s *Server) Start(acceptwg *sync.WaitGroup) error {
	if acceptwg != nil {
		defer acceptwg.Done()
		s.acceptwg = acceptwg
	}
	err := s.listener.Accept(handleMsg, handleDisconn)
	if err != nil && s.debug {
		log.Fatalln("Listener returned an error while closing:", err)
	}

	s.mutex.Lock()
	s.err = err
	s.mutex.Unlock()

	return err
}

// Stop stops a server instance gracefully. For listener is closed to deny any new connections, then server waits for all connections to be
// closed gracefully to deny any new requests, and finally it waits for all pending requests to be finished.
func (s *Server) Stop() error {
	// close the listener and wait for listener.Accept to return
	err := s.listener.Close()
	if s.acceptwg != nil {
		s.acceptwg.Wait()
	}

	// close all active connections discarding any read/writes that is going on currently. this is not a problem as we always require an ACK
	for _, user := range users {
		err := user.Conn.Close()
		if err != nil {
			return err
		}
		user.Conn = nil
	}
	for _, conn := range s.listener.Conns {
		err := conn.Close()
		if err != nil {
			return err
		}
	}
	s.connwg.Wait()

	// wait for all pending requests to be handled/finalized
	s.reqwg.Wait()

	if s.err != nil {
		return s.err
	}
	return err
}

// handleMsg handles incoming client messages.
func handleMsg(conn *Conn, session *Session, msg []byte) {
	// authenticate the session if not already done
	if session.UserID == 0 {
		userID, err := auth(conn.ConnectionState().PeerCertificates, msg)
		if err != nil {
			session.Error = fmt.Errorf("Cannot parse client message or method mismatched: %v", err)
		}
		session.UserID = userID
		users[userID].Conn = conn
		// todo: ack auth message, start sending other queued messages one by one
		// can have 2 approaches here
		// 1. users[id].send(...) & users[id].queue(...)
		// 2. conn.write(...) && queue[id].conn = ...
		return
	}

	// queue the incoming request and send an ack
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

func handleDisconn(conn *Conn, session *Session) {
	users[session.UserID].Conn = nil
}
