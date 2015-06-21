package devastator

import (
	"log"
	"sync"

	"github.com/nbusy/neptulon"
	"github.com/nbusy/neptulon/jsonrpc"
)

var users = make(map[uint32]*User)

// Server wraps a listener instance and registers default connection and message handlers with the listener.
type Server struct {
	debug    bool
	err      error
	neptulon *neptulon.App
	mutex    sync.Mutex
}

// NewServer creates and returns a new server instance with a listener created using given parameters.
// Debug mode dumps raw TCP data to stderr (log.Println() default).
func NewServer(cert, privKey []byte, laddr string, debug bool) (*Server, error) {
	nep, err := neptulon.NewApp(cert, privKey, laddr, debug)
	if err != nil {
		return nil, err
	}

	jrpc, err := jsonrpc.NewApp(nep)
	if err != nil {
		return nil, err
	}

	pubrout, err := jsonrpc.NewRouter(jrpc)
	if err != nil {
		return nil, err
	}

	pubrout.Request("close", func(ctx *jsonrpc.ReqContext) {
		ctx.Res = "ACK" // should be notification and should close conn immediately
	})
	pubrout.Request("auth.cert", func(ctx *jsonrpc.ReqContext) {
		certs := ctx.Conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			ctx.ResErr = &jsonrpc.ResError{Code: 666, Message: "Invalid client certificate.", Data: certs}
		} else {
			ctx.Res = "OK"
		}
	})

	_, err = jsonrpc.NewCertAuth(jrpc)
	if err != nil {
		return nil, err
	}

	privrout, err := jsonrpc.NewRouter(jrpc)
	if err != nil {
		return nil, err
	}
	privrout.Request("echo", func(ctx *jsonrpc.ReqContext) {
		ctx.Res = ctx.Req.Params
	})

	// n.Middleware() // json rpc protocol
	// p.Middleware("auth.google") // public json rpc routes
	// p.Middleware() // cert auth
	// p.Middleware() // private json rpc routes
	// p.Middleware() // 404-like handler
	// p.Middleware() // request-response logger

	return &Server{
		debug:    debug,
		neptulon: nep,
	}, nil
}

// Start starts accepting connections on the internal listener and handles connections with registered onnection and message handlers.
// This function blocks and never returns, unless there was an error while accepting a new connection or the server was closed.
func (s *Server) Start() error {
	err := s.neptulon.Run()
	if err != nil && s.debug {
		log.Fatalln("Listener returned an error while closing:", err)
	}

	s.mutex.Lock()
	s.err = err
	s.mutex.Unlock()

	return err
}

// Stop stops a server instance.
func (s *Server) Stop() error {
	err := s.neptulon.Stop()

	// close all active connections discarding any read/writes that is going on currently
	// this is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server
	for _, user := range users {
		err := user.Conn.Close()
		if err != nil {
			return err
		}
		user.Conn = nil
	}

	s.mutex.Lock()
	if s.err != nil {
		return s.err
	}
	s.mutex.Unlock()
	return err
}

// ------------------------- legacy ----------------------------

// // handleMsg handles incoming client messages.
// func handleMsg(conn *neptulon.Conn, msg []byte) {
// 	// authenticate the session if not already done
// 	if session.UserID == 0 {
// 		userID, err := auth(conn.ConnectionState().PeerCertificates, msg)
// 		if err != nil {
// 			session.Error = fmt.Errorf("Cannot parse client message or method mismatched: %v", err)
// 		}
// 		session.UserID = userID
// 		users[userID].Conn = conn
// 		// todo: ack auth message, start sending other queued messages one by one
// 		// can have 2 approaches here
// 		// 1. users[id].send(...) & users[id].queue(...)
// 		// 2. conn.write(...) && queue[id].conn = ...
// 		return
// 	}
//
// 	// router: process the message and queue a reply if necessary and send an ack
// }
//
// // auth handles Google+ sign-in and client certificate authentication.
// func auth(peerCerts []*x509.Certificate, msg []byte) (userID uint32, err error) {
// 	// client certificate authorization: certificate is verified by the TLS listener instance so we trust it
// 	if len(peerCerts) > 0 {
// 		idstr := peerCerts[0].Subject.CommonName
// 		uid64, err := strconv.ParseUint(idstr, 10, 32)
// 		if err != nil {
// 			return 0, err
// 		}
// 		userID = uint32(uid64)
// 		log.Printf("Client connected with client certificate subject: %+v", peerCerts[0].Subject)
// 		return userID, nil
// 	}
//
// 	// Google+ authentication
// 	var req jsonrpc.Request
// 	if err = json.Unmarshal(msg, &req); err != nil {
// 		return
// 	}
//
// 	switch req.Method {
// 	case "auth.token":
// 		var token string
// 		if err = json.Unmarshal(req.Params, &token); err != nil {
// 			return
// 		}
// 		// assume that token = user ID for testing
// 		uid64, err := strconv.ParseUint(token, 10, 32)
// 		if err != nil {
// 			return 0, err
// 		}
// 		userID = uint32(uid64)
// 		return userID, nil
// 	case "auth.google":
// 		// todo: ping google, get user info, save user info in DB, generate and return permanent jwt token (or should this part be NBusy's business?)
// 		return
// 	default:
// 		return 0, errors.New("initial unauthenticated request should be in the 'auth.xxx' form")
// 	}
// }
//
// func handleDisconn(conn *neptulon.Conn) {
// 	users[session.UserID].Conn = nil
// }
