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
	// neptulon framework components
	neptulon  *neptulon.Server
	jsonrpc   *jsonrpc.Server
	pubRoute  *jsonrpc.Router
	privRoute *jsonrpc.Router

	// devastator server components
	db      DB
	certMgr CertMgr
	queue   Queue

	debug    bool  // dump raw TCP message to stderr using log.Println()
	err      error // last error returned by neptulon framework before closing listener
	errMutex sync.Mutex
}

// NewServer creates and returns a new server instance with a listener created using given parameters.
// Debug mode dumps raw TCP messages to stderr using log.Println().
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
	}

	s.jsonrpc, err = jsonrpc.NewServer(nep)
	if err != nil {
		return nil, err
	}

	s.pubRoute, err = jsonrpc.NewRouter(s.jsonrpc)
	if err != nil {
		return nil, err
	}

	initPubRoutes(s.pubRoute, s.db, &s.certMgr)

	// --- all requests below this point must be authenticated ---

	CertAuth(s.jsonrpc)

	s.queue = NewQueue(&s) // todo: doing this like this is really weird (queue middleware can be separated from queue type)

	s.privRoute, err = jsonrpc.NewRouter(s.jsonrpc)
	if err != nil {
		return nil, err
	}

	initPrivRoutes(s.privRoute, &s.queue)

	nep.Disconn(func(c *neptulon.Conn) {
		if id, ok := c.Data.GetOk("userid"); ok {
			s.queue.RemoveConn(id.(string))
		}
	})

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
