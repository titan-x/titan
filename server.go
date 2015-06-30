package devastator

import (
	"io/ioutil"
	"log"
	"net/http"
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
// Debug mode dumps raw TCP data to stderr using log.Println().
func NewServer(cert, privKey []byte, laddr string, debug bool) (*Server, error) {
	nep, err := neptulon.NewApp(cert, privKey, laddr, debug)
	if err != nil {
		return nil, err
	}

	rpc, err := jsonrpc.NewApp(nep)
	if err != nil {
		return nil, err
	}

	pubrout, err := jsonrpc.NewRouter(rpc)
	if err != nil {
		return nil, err
	}

	// retrieve user info (display name, e-mail, profile pic) using an access token that has 'profile' and 'email' scopes
	pubrout.Request("auth.google", func(ctx *jsonrpc.ReqContext) {
		res, err := http.Get("https://www.googleapis.com/plus/v1/people/me?access_token=" + ctx.Req.Params.(map[string]string)["token"])
		if err != nil {
			log.Fatal(err)
		}

		profile, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s", profile)

		// email: profile.emails[0].value,
		// name: profile.displayName,
		// picture: (yield request.get(profile.image.url, {encoding: 'base64'})).body

		// if authenticated generate "userid", set it in session, create and send client-certificate as reponse
	})

	_, err = jsonrpc.NewCertAuth(rpc)
	if err != nil {
		return nil, err
	}

	privrout, err := jsonrpc.NewRouter(rpc)
	if err != nil {
		return nil, err
	}

	privrout.Request("echo", func(ctx *jsonrpc.ReqContext) {
		ctx.Res = ctx.Req.Params
		if ctx.Res == nil {
			ctx.Res = ""
		}
	})

	// p.Middleware(NotFoundHandler()) // 404-like handler
	// p.Middleware(Logger()) // request-response logger (the pointer fields in request/response objects will have to change for this to work)

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
