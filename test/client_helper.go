package test

import (
	"testing"

	"github.com/nb-titan/titan"
	"github.com/neptulon/jsonrpc/test"
)

// ClientHelper is a Titan Client wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ClientHelper struct {
	Client *titan.Client

	jrpcCH        *test.ClientHelper // inner Neptulon JSON-RPC ClientHelper object
	testing       *testing.T
	ca, cert, key []byte
}

// NewClientHelper creates a new client helper object.
func NewClientHelper(t *testing.T, ca []byte, addr string) *ClientHelper {
	jrpcCH := test.NewClientHelper(t, addr)
	c := titan.UseClient(jrpcCH.Client)
	return &ClientHelper{
		Client:  c,
		jrpcCH:  jrpcCH,
		testing: t,
		ca:      ca,
	}
}

// Connect connects to a server.
func (ch *ClientHelper) Connect() *ClientHelper {
	ch.Client.UseTLS(ch.ca, ch.cert, ch.key)
	ch.jrpcCH.Connect()
	return ch
}

// AsUser attaches given user's client certificate and private key to the connection.
func (ch *ClientHelper) AsUser(u *titan.User) *ClientHelper {
	ch.cert = u.Cert
	ch.key = u.Key
	return ch
}

// Close closes a connection.
func (ch *ClientHelper) Close() {
	ch.jrpcCH.Close()
}
