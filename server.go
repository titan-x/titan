package titan

import (
	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
	"github.com/neptulon/neptulon/middleware/jwt"
	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/data/inmem"
)

// Server wraps a listener instance and registers default connection and message handlers with the listener.
type Server struct {
	// neptulon framework components
	neptulon   *neptulon.Server
	pubRouter  *middleware.Router
	privRouter *middleware.Router

	// titan server components
	db    data.DB
	queue data.Queue
}

// NewServer creates a new server.
func NewServer(addr string) (*Server, error) {
	if (Conf == Config{}) {
		InitConf("")
	}

	s := Server{neptulon: neptulon.NewServer(addr)}

	if err := s.SetDB(inmem.NewDB()); err != nil {
		return nil, err
	}
	if err := s.SetQueue(inmem.NewQueue(s.neptulon.SendRequest)); err != nil {
		return nil, err
	}

	s.neptulon.MiddlewareFunc(middleware.Logger)
	s.pubRouter = middleware.NewRouter()
	s.neptulon.Middleware(s.pubRouter)
	initPubRoutes(s.pubRouter, &s.db, Conf.App.JWTPass())

	//all communication below this point is authenticated
	s.neptulon.MiddlewareFunc(jwt.HMAC(Conf.App.JWTPass()))
	s.neptulon.Middleware(s.queue)
	s.privRouter = middleware.NewRouter()
	s.neptulon.Middleware(s.privRouter)
	initPrivRoutes(s.privRouter, &s.queue)
	// todo: r.Middleware(NotFoundHandler()) - 404-like handler, if any request reaches this point without being handled

	s.neptulon.DisconnHandler(func(c *neptulon.Conn) {
		// only handle this event for previously authenticated
		if id, ok := c.Session.GetOk("userid"); ok {
			s.queue.RemoveConn(id.(string))
		}
	})

	return &s, nil
}

// SetDB sets the database implementation to be used by the server. If not supplied, in-memory database implementation is used.
func (s *Server) SetDB(db data.DB) error {
	if err := db.Seed(false, Conf.App.JWTPass()); err != nil {
		return err
	}

	s.db = db
	return nil
}

// SetQueue sets the queue implementation to be used by the server. If not supplied, in-memory queue implementation is used.
func (s *Server) SetQueue(queue data.Queue) error {
	s.queue = queue
	return nil
}

// ListenAndServe starts the Titan server. This function blocks until server is closed.
func (s *Server) ListenAndServe() error {
	return s.neptulon.ListenAndServe()
}

// Close the server and all of the active connections, discarding any read/writes that is going on currently.
// This is not a problem as we always require an ACK but it will also mean that message deliveries will be at-least-once; to-and-from the server.
func (s *Server) Close() error {
	return s.neptulon.Close()
}
