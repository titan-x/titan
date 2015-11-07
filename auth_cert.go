package titan

import (
	"log"

	"github.com/neptulon/jsonrpc"
	"github.com/neptulon/neptulon"
)

// CertAuth does TLS client-auth check and sets user ID in connection session store.
// Connection is closed without returning any reason if cert is invalid.
func CertAuth(server *jsonrpc.Server) {
	server.ReqMiddleware(func(ctx *jsonrpc.ReqCtx) {
		if !certAuth(ctx.Conn) {
			ctx.Done = true
		}
	})
	server.ResMiddleware(func(ctx *jsonrpc.ResCtx) {
		if !certAuth(ctx.Conn) {
			ctx.Done = true
		}
	})
	server.NotMiddleware(func(ctx *jsonrpc.NotCtx) {
		if !certAuth(ctx.Conn) {
			ctx.Done = true
		}
	})
}

// certAuth does client-auth check and sets user ID in connection session store.
// Returns true if authentication is successful.
func certAuth(c *neptulon.Conn) bool {
	if _, ok := c.Data.GetOk("userid"); ok {
		return true
	}

	// if provided, client certificate is verified by the TLS listener so the peerCerts list in the connection is trusted
	certs := c.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		log.Println("Invalid client-certificate authentication attempt:", c.RemoteAddr())
		c.Close()
		return false
	}

	userID := certs[0].Subject.CommonName
	c.Data.Set("userid", userID)
	log.Printf("Client authenticated. TLS/IP: %v, User ID: %v, Conn ID: %v\n", c.RemoteAddr(), userID, c.ID)
	return true
}
