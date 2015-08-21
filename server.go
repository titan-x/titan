package devastator

import (
	"fmt"
	"log"
	"sync"

	"github.com/nbusy/neptulon"
	"github.com/nbusy/neptulon/jsonrpc"
)

// Server wraps a listener instance and registers default connection and message handlers with the listener.
type Server struct {
	debug    bool
	err      error
	neptulon *neptulon.App
	mutex    sync.Mutex
	db       DB
	certMgr  *CertMgr
}

// NewServer creates and returns a new server instance with a listener created using given parameters.
// Debug mode dumps raw TCP data to stderr using log.Println().
func NewServer(cert, privKey, clientCACert, clientCAKey []byte, laddr string, debug bool) (*Server, error) {
	nep, err := neptulon.NewApp(cert, privKey, clientCACert, laddr, debug)
	if err != nil {
		return nil, err
	}

	s := Server{
		debug:    debug,
		neptulon: nep,
		db:       NewInMemDB(),
		certMgr:  NewCertMgr(clientCACert, clientCAKey),
	}

	rpc, err := jsonrpc.NewApp(nep)
	if err != nil {
		return nil, err
	}

	pubRoute, err := jsonrpc.NewRouter(rpc)
	if err != nil {
		return nil, err
	}

	pubRoute.Request("auth.google", func(ctx *jsonrpc.ReqContext) {
		googleAuth(ctx, s.db, s.certMgr)
	})

	pubRoute.Notification("close", func(ctx *jsonrpc.NotContext) {
		ctx.Done = true
		ctx.Conn.Close()
	})

	// pubRoute.NotFound(...)
	// todo: if the first incoming message in public route is not one of close/google.auth,
	// close the connection right away (and maybe wait for client to return ACK then close?)

	_, err = jsonrpc.NewCertAuth(rpc)
	if err != nil {
		return nil, err
	}

	privRoute, err := jsonrpc.NewRouter(rpc)
	if err != nil {
		return nil, err
	}

	privRoute.Request("echo", func(ctx *jsonrpc.ReqContext) {
		ctx.Res = ctx.Req.Params
		if ctx.Res == nil {
			ctx.Res = ""
		}
	})

	// privRoute.Middleware(NotFoundHandler()) // 404-like handler
	// privRoute/pubRoute.Middleware(Logger()) // request-response logger (the pointer fields in request/response objects will have to change for this to work)

	return &s, nil
}

// UseDB sets the database to be used by the server. If not supplied, in-memory database implementation is used.
func (s *Server) UseDB(db DB) error {
	s.db = db
	return nil
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

// Stop stops the server and closes all of the active connections discarding any read/writes that is going on currently.
// This is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server.
func (s *Server) Stop() error {
	err := s.neptulon.Stop()

	s.mutex.Lock()
	if s.err != nil {
		return fmt.Errorf("Past internal error: %v", s.err)
	}
	s.mutex.Unlock()
	return err
}
