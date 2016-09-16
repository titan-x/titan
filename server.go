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
	pubRoutes  *middleware.Router
	privRoutes *middleware.Router

	// titan server components
	db    data.DB
	queue data.Queue
}

// NewServer creates a new server.
func NewServer(addr string) (*Server, error) {
	if (Conf == Config{}) {
		InitConf("")
	}

	n := neptulon.NewServer(addr)
	s := Server{
		neptulon: n,
		queue:    inmem.NewQueue(n.SendRequest),
	}

	if err := s.SetDB(inmem.NewDB()); err != nil {
		return nil, err
	}

	s.neptulon.MiddlewareFunc(middleware.Logger)
	s.pubRoutes = middleware.NewRouter()
	s.neptulon.Middleware(s.pubRoutes)
	initPubRoutes(s.pubRoutes, &s.db, Conf.App.JWTPass())

	//all communication below this point is authenticated
	s.neptulon.MiddlewareFunc(jwt.HMAC(Conf.App.JWTPass()))
	s.neptulon.Middleware(s.queue)
	s.privRoutes = middleware.NewRouter()
	s.neptulon.Middleware(s.privRoutes)
	initPrivRoutes(s.privRoutes, s.queue)
	// r.Middleware(NotFoundHandler()) - 404-like handler

	s.neptulon.DisconnHandler(func(c *neptulon.Conn) {
		// only handle this event for previously authenticated
		if id, ok := c.Session.GetOk("userid"); ok {
			s.queue.RemoveConn(id.(string))
		}
	})

	return &s, nil
}

// SetDB sets the database to be used by the server. If not supplied, in-memory database implementation is used.
func (s *Server) SetDB(db data.DB) error {
	if err := db.Seed(false, Conf.App.JWTPass()); err != nil {
		return err
	}

	s.db = db
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
