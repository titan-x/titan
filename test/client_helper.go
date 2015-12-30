package test

import (
	"io"

	"github.com/nb-titan/titan"
	"github.com/neptulon/jsonrpc"
)

// ConnHelper is a neptulon/jsonrpc.Conn wrapper.
// All the functions are wrapped with proper test runner error logging.
type ConnHelper struct {
	conn      *jsonrpc.Client
	cert, key []byte
}

// AsUser attaches given user's client certificate and private key to the connection.
func (c *ConnHelper) AsUser(u *titan.User) *ConnHelper {
	return c.WithCert(u.Cert, u.Key)
}

// WithCert attaches given PEM encoded client certificate and private key to the connection.
func (c *ConnHelper) WithCert(cert, key []byte) *ConnHelper {
	c.cert = cert
	c.key = key
	return c
}

// VerifyConnClosed verifies that the connection is in closed state.
// Verification is done via reading from the channel and checking that returned error is io.EOF.
func (c *ConnHelper) VerifyConnClosed() bool {
	_, _, _, err := c.conn.ReadMsg(nil, nil)
	if err != io.EOF {
		return false
	}

	return true
}
