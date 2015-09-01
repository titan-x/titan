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
	neptulon  *neptulon.Server
	jsonrpc   *jsonrpc.Server
	pubRoute  *jsonrpc.Router
	privRoute *jsonrpc.Router

	db      DB
	certMgr *CertMgr
	queue   Queue

	debug    bool
	err      error
	errMutex sync.Mutex
}

// NewServer creates and returns a new server instance with a listener created using given parameters.
// Debug mode dumps raw TCP data to stderr using log.Println().
func NewServer(cert, privKey, clientCACert, clientCAKey []byte, laddr string, debug bool) (*Server, error) {
	nep, err := neptulon.NewServer(cert, privKey, clientCACert, laddr, debug)
	if err != nil {
		return nil, err
	}

	s := Server{
		debug:    debug,
		neptulon: nep,
		db:       NewInMemDB(),
		certMgr:  NewCertMgr(clientCACert, clientCAKey),
		queue:    NewQueue(),
	}

	s.jsonrpc, err = jsonrpc.NewServer(nep)
	if err != nil {
		return nil, err
	}

	s.pubRoute, err = jsonrpc.NewRouter(s.jsonrpc)
	if err != nil {
		return nil, err
	}

	s.pubRoute.Request("auth.google", func(ctx *jsonrpc.ReqCtx) {
		googleAuth(ctx, s.db, s.certMgr)
	})

	s.pubRoute.Notification("conn.close", func(ctx *jsonrpc.NotCtx) {
		ctx.Done = true
		ctx.Conn.Close()
	})

	// pubRoute.NotFound(...)
	// todo: if the first incoming message in public route is not one of close/google.auth,
	// close the connection right away (and maybe wait for client to return ACK then close?)

	_, err = NewCertAuth(s.jsonrpc)
	if err != nil {
		return nil, err
	}

	s.privRoute, err = jsonrpc.NewRouter(s.jsonrpc)
	if err != nil {
		return nil, err
	}

	s.privRoute.Request("msg.echo", func(ctx *jsonrpc.ReqCtx) {
		ctx.Params(&ctx.Res)
		if ctx.Res == nil {
			ctx.Res = ""
		}
	})

	type msgSendRequest struct {
		to      string
		message string
	}

	s.privRoute.Request("msg.send", func(ctx *jsonrpc.ReqCtx) {
		// try to send the incoming message right away
		var r msgSendRequest
		ctx.Params(&r)
		s.queue.AddRequest(r.to, &jsonrpc.Request{ID: "456", Method: "msg.recv", Params: r.message})
	})

	s.privRoute.Request("msg.recv", func(ctx *jsonrpc.ReqCtx) {

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

	s.errMutex.Lock()
	s.err = err
	s.errMutex.Unlock()

	return err
}

// Stop stops the server and closes all of the active connections discarding any read/writes that is going on currently.
// This is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server.
func (s *Server) Stop() error {
	err := s.neptulon.Stop()

	s.errMutex.Lock()
	if s.err != nil {
		return fmt.Errorf("Past internal error: %v", s.err)
	}
	s.errMutex.Unlock()
	return err
}
