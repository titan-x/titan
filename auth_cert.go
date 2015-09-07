package devastator

import (
	"log"

	"github.com/nbusy/neptulon/jsonrpc"
)

// todo: action taken in response to un authenticated req/res/not messages should be configurable on per-item basis,
// as closing the connection always might not be the desired behavior

// todo2: we need to pass client CA cert as a param here which will add it to listener.tls.Config file as client CA cert
// rather than TLS listener always requiring client CA cert w/ constructor

// CertAuth is a TLS certificate authentication middleware for Neptulon JSON-RPC server.
type CertAuth struct {
}

// NewCertAuth creates and registers a new certificate authentication middleware instance with a Neptulon JSON-RPC server.
func NewCertAuth(server *jsonrpc.Server) (*CertAuth, error) {
	a := CertAuth{}
	server.ReqMiddleware(a.reqMiddleware)
	server.ResMiddleware(a.resMiddleware)
	server.NotMiddleware(a.notMiddleware)
	return &a, nil
}

func (a *CertAuth) reqMiddleware(ctx *jsonrpc.ReqCtx) {
	if _, ok := ctx.Conn.Data.Get("userid"); ok {
		return
	}

	// if provided, client certificate is verified by the TLS listener so the peerCerts list in the connection is trusted
	certs := ctx.Conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		log.Println("Invalid client-certificate authentication attempt:", ctx.Conn.RemoteAddr())
		ctx.Done = true
		ctx.Conn.Close()
		return
	}

	userID := certs[0].Subject.CommonName
	ctx.Conn.Data.Set("userid", userID)
	log.Println("Client-certificate authenticated:", ctx.Conn.RemoteAddr(), userID)
}

func (a *CertAuth) resMiddleware(ctx *jsonrpc.ResCtx) {
	if _, ok := ctx.Conn.Data.Get("userid"); ok {
		return
	}

	ctx.Done = true
	ctx.Conn.Close()
}

func (a *CertAuth) notMiddleware(ctx *jsonrpc.NotCtx) {
	if _, ok := ctx.Conn.Data.Get("userid"); ok {
		return
	}

	ctx.Done = true
	ctx.Conn.Close()
}
