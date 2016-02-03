package titan

import (
	"fmt"
	"log"
	"sync"

	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
)

// Server wraps a listener instance and registers default connection and message handlers with the listener.
type Server struct {
	// neptulon framework components
	server    *neptulon.Server
	pubRoute  *middleware.Router
	privRoute *middleware.Router

	// titan server components
	db    DB
	queue Queue

	debug    bool  // dump raw TCP message to stderr using log.Println()
	err      error // last error returned by neptulon framework before closing listener
	errMutex sync.Mutex
}

// NewServer creates and returns a new server instance with a listener created using given parameters.
// Debug mode dumps raw TCP messages to stderr using log.Println().
func NewServer(addr string) (*Server, error) {
	s := Server{
		server: neptulon.NewServer(addr),
		db:     NewInMemDB(),
		queue:  NewQueue(),
	}

	s.pubRoute = middleware.NewRouter()
	s.server.Middleware(s.pubRoute)
	initPubRoutes(s.pubRoute, s.db)

	// --- all communication below this point is authenticated --- //

	// CertAuth(s.neptulon)
	s.queue.Middleware(s.server)
	s.privRoute = middleware.NewRouter()
	s.server.Middleware(s.privRoute)
	initPrivRoutes(s.privRoute, &s.queue)

	s.queue.SetServer(s.server) // todo: research a better way to handle inner-circular dependencies so remove these lines back into Server contructor (maybe via dereferencing: http://openmymind.net/Things-I-Wish-Someone-Had-Told-Me-About-Go/, but then initializers actually using the pointer values would have to be lazy!)

	s.server.DisconnHandler(func(c *neptulon.Conn) {
		// only handle this event for previously authenticated
		if id, ok := c.Session.GetOk("userid"); ok {
			s.queue.RemoveConn(id.(string))
		}
	})

	return &s, nil
}

// SetDB sets the database to be used by the server. If not supplied, in-memory database implementation is used.
func (s *Server) SetDB(db DB) error {
	s.db = db
	return nil
}

// Start starts accepting connections on the internal listener and handles connections with registered onnection and message handlers.
// This function blocks and never returns, unless there was an error while accepting a new connection or the server was closed.
func (s *Server) Start() error {
	err := s.server.Start()
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
	err := s.server.Close()

	s.errMutex.Lock()
	if s.err != nil {
		return fmt.Errorf("Past internal error: %v", s.err)
	}
	s.errMutex.Unlock()
	return err
}
