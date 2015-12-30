package titan

import (
	"log"

	"github.com/neptulon/jsonrpc"
)

// CertAuth does TLS client-auth check and sets user ID in connection session store.
// Connection is closed without returning any reason if cert is invalid.
func CertAuth(server *jsonrpc.Server) {
	server.ReqMiddleware(func(ctx *jsonrpc.ReqCtx) error {
		if !certAuth(ctx.Client) {
			return nil
		}
		return ctx.Next()
	})
	server.ResMiddleware(func(ctx *jsonrpc.ResCtx) error {
		if !certAuth(ctx.Client) {
			return nil
		}
		return ctx.Next()
	})
	server.NotMiddleware(func(ctx *jsonrpc.NotCtx) error {
		if !certAuth(ctx.Client) {
			return nil
		}
		return ctx.Next()
	})
}

// certAuth does client-auth check and sets user ID in connection session store.
// Returns true if authentication is successful.
func certAuth(c *jsonrpc.Client) bool {
	if _, ok := c.Session().GetOk("userid"); ok {
		return true
	}

	// if provided, client certificate is verified by the TLS listener so the peerCerts list in the connection is trusted
	connState, _ := c.Conn.ConnectionState()
	certs := connState.PeerCertificates
	if len(certs) == 0 {
		log.Println("Invalid client-certificate authentication attempt:", c.Conn.RemoteAddr())
		c.Close()
		return false
	}

	userID := certs[0].Subject.CommonName
	c.Session().Set("userid", userID)
	log.Printf("Client authenticated. TLS/IP: %v, User ID: %v, Conn ID: %v\n", c.Conn.RemoteAddr(), userID, c.ConnID())
	return true
}
