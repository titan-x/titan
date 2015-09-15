package test

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/nbusy/devastator"
	"github.com/nbusy/neptulon/jsonrpc"
)

// ClientHelper is a neptulon/jsonrpc.Client wrapper.
// All the functions are wrapped with proper test runner error logging.
type ClientHelper struct {
	client    *jsonrpc.Client
	server    *ServerHelper
	testing   *testing.T
	cert, key []byte
}

// NewClientHelper creates a new client helper object.
func NewClientHelper(t *testing.T, s *ServerHelper) *ClientHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	return &ClientHelper{testing: t, server: s}
}

// AsUser attaches given user's client certificate and private key to the connection.
func (c *ClientHelper) AsUser(u *devastator.User) *ClientHelper {
	c.cert = u.Cert
	c.key = u.Key
	return c
}

// WithCert attaches given PEM encoded client certificate and private key to the connection.
func (c *ClientHelper) WithCert(cert, key []byte) *ClientHelper {
	c.cert = cert
	c.key = key
	return c
}

// Dial initiates a connection.
func (c *ClientHelper) Dial() *ClientHelper {
	addr := "127.0.0.1:" + devastator.Conf.App.Port

	// retry connect in case we're operating on a very slow machine
	for i := 0; i <= 5; i++ {
		client, err := jsonrpc.Dial(addr, c.server.IntCACert, c.cert, c.key, false) // no need for debug mode on client conn as we have it on server conn already
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

// WriteRequest sends a request message through the client connection.
func (c *ClientHelper) WriteRequest(method string, params interface{}) (reqID string) {
	id, err := c.client.WriteRequest(method, params)
	if err != nil {
		c.testing.Fatal("Failed to write request to client connection:", err)
	}
	return id
}

// WriteRequestArr sends a request message through the client connection. Params object is variadic.
func (c *ClientHelper) WriteRequestArr(method string, params ...interface{}) (reqID string) {
	return c.WriteRequest(method, params)
}

// WriteNotification sends a notification message through the client connection.
func (c *ClientHelper) WriteNotification(method string, params interface{}) {
	if err := c.client.WriteNotification(method, params); err != nil {
		c.testing.Fatal("Failed to write notification to client connection:", err)
	}
}

// WriteNotificationArr sends a notification message through the client connection. Params object is variadic.
func (c *ClientHelper) WriteNotificationArr(method string, params ...interface{}) {
	c.WriteNotification(method, params)
}

// WriteResponse sends a response message through the client connection.
func (c *ClientHelper) WriteResponse(id string, result interface{}, err *jsonrpc.ResError) {
	if err := c.client.WriteResponse(id, result, err); err != nil {
		c.testing.Fatal("Failed to write response to client connection:", err)
	}
}

// ReadReq reads a request object from a client connection.
// If incoming message is not a request, a fatal error is logged.
// Optionally, you can pass in a data structure that the returned JSON-RPC request params data will be serialized into.
// Otherwise json.Unmarshal defaults apply.
func (c *ClientHelper) ReadReq(paramsData interface{}) *jsonrpc.Request {
	req, _, _, err := c.client.ReadMsg(nil, paramsData)
	if err != nil {
		c.testing.Fatal("Failed to read request from client connection:", err)
	}

	return req
}

// ReadRes reads a response object from a client connection.
// If incoming message is not a response, a fatal error is logged.
// Optionally, you can pass in a data structure that the returned JSON-RPC response result data will be serialized into.
// Otherwise json.Unmarshal defaults apply.
func (c *ClientHelper) ReadRes(resultData interface{}) *jsonrpc.Response {
	_, res, _, err := c.client.ReadMsg(resultData, nil)
	if err != nil {
		c.testing.Fatal("Failed to read response from client connection:", err)
	}

	return res
}

// ReadMsg reads a JSON-RPC message from a client connection. If possible, use more specific ReadReq/ReadRes/ReadNot methods instead.
// Optionally, you can pass in a data structure that the returned JSON-RPC response result data will be serialized into (same for request params).
// Otherwise json.Unmarshal defaults apply.
func (c *ClientHelper) ReadMsg(resultData interface{}, paramsData interface{}) (req *jsonrpc.Request, res *jsonrpc.Response, not *jsonrpc.Notification) {
	req, res, not, err := c.client.ReadMsg(resultData, paramsData)
	if err != nil {
		c.testing.Fatal("Failed to read message from client connection:", err)
	}

	return
}

// VerifyConnClosed verifies that the connection is in closed state.
// Verification is done via reading from the channel and checking that returned error is io.EOF.
func (c *ClientHelper) VerifyConnClosed() bool {
	_, _, _, err := c.client.ReadMsg(nil, nil)
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
