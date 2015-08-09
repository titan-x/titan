package test

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/nbusy/devastator"
	"github.com/nbusy/neptulon/jsonrpc"
)

// ClientHelper is a JSON-RPC client wrapper with built-in error logging for testing.
type ClientHelper struct {
	client    *jsonrpc.Client
	testing   *testing.T
	cert, key []byte
}

// NewClientHelper creates a new JSON-RPC client wrapper which has built-in error logging for testing.
func NewClientHelper(t *testing.T) *ClientHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	return &ClientHelper{testing: t}
}

// DefaultCert attaches default test client certificate to the connection.
func (c *ClientHelper) DefaultCert() *ClientHelper {
	c.cert = certChain.ClientCert
	c.key = certChain.ClientKey
	return c
}

// Cert attaches given PEM encoded client certificate to the connection.
func (c *ClientHelper) Cert(cert, key []byte) *ClientHelper {
	c.cert = cert
	c.key = key
	return c
}

// Dial initiates a connection.
func (c *ClientHelper) Dial() *ClientHelper {
	addr := "127.0.0.1:" + devastator.Conf.App.Port

	// retry connect in case we're operating on a very slow machine
	for i := 0; i <= 5; i++ {
		client, err := jsonrpc.Dial(addr, certChain.IntCAKey, c.cert, c.key, false) // no need for debug mode on client conn
		if err != nil {
			if operr, ok := err.(*net.OpError); ok && operr.Op == "dial" && operr.Err.Error() == "connection refused" {
				time.Sleep(time.Millisecond * 50)
				continue
			} else if i == 5 {
				c.testing.Fatalf("Cannot connect to server address %v after 5 retries, with error: %v", addr, err)
			}
			c.testing.Fatalf("Cannot connect to server address %v with error: %v", addr, err)
		}

		if i != 0 {
			c.testing.Logf("WARNING: it took %v retries to connect to the server, which might indicate code issues or slow machine.", i)
		}

		client.SetReadDeadline(10)
		c.client = client
		return c
	}

	return nil
}

// WriteRequest writes a request to a client connection with error logging for testing.
func (c *ClientHelper) WriteRequest(method string, params interface{}) (reqID string) {
	id, err := c.client.WriteRequest(method, params)
	if err != nil {
		c.testing.Fatal("Failed to write request to client connection:", err)
	}
	return id
}

// ReadMsg reads a JSON-RPC message from a client connection with error logging for testing.
func (c *ClientHelper) ReadMsg() (req *jsonrpc.Request, res *jsonrpc.Response, not *jsonrpc.Notification) {
	req, res, not, err := c.client.ReadMsg()
	if err != nil {
		c.testing.Fatal("Failed to read message from client connection:", err)
	}

	return
}

// ReadRes reads a response object from a client connection. If incoming message is not a response, an error is logged.
func (c *ClientHelper) ReadRes() *jsonrpc.Response {
	_, res, _, err := c.client.ReadMsg()
	if err != nil {
		c.testing.Fatal("Failed to read response from client connection:", err)
	}

	return res
}

// VerifyConnClosed verifies that the connection is in closed state.
// Verification is done via reading from the channel and checking that returned error is io.EOF.
func (c *ClientHelper) VerifyConnClosed() bool {
	_, _, _, err := c.client.ReadMsg()
	if err != io.EOF {
		return false
	}

	return true
}

// Close closes a client connection.
func (c *ClientHelper) Close() {
	if err := c.client.Close(); err != nil {
		c.testing.Fatal("Failed to close client connection:", err)
	}
}
