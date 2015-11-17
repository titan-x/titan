package test

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/nb-titan/titan"
	"github.com/neptulon/jsonrpc"
)

// ConnHelper is a neptulon/jsonrpc.Conn wrapper.
// All the functions are wrapped with proper test runner error logging.
type ConnHelper struct {
	conn      *jsonrpc.Conn
	server    *ServerHelper // server that this connection will be made to
	testing   *testing.T
	cert, key []byte
}

// NewConnHelper creates a new connection helper object.
// Takes target server as an argument to retrieve server certs, address, etc.
func NewConnHelper(t *testing.T, s *ServerHelper) *ConnHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	return &ConnHelper{testing: t, server: s}
}

// AsUser attaches given user's client certificate and private key to the connection.
func (c *ConnHelper) AsUser(u *titan.User) *ConnHelper {
	c.cert = u.Cert
	c.key = u.Key
	return c
}

// WithCert attaches given PEM encoded client certificate and private key to the connection.
func (c *ConnHelper) WithCert(cert, key []byte) *ConnHelper {
	c.cert = cert
	c.key = key
	return c
}

// Dial initiates a connection.
func (c *ConnHelper) Dial() *ConnHelper {
	addr := "127.0.0.1:" + titan.Conf.App.Port

	// retry connect in case we're operating on a very slow machine
	for i := 0; i <= 5; i++ {
		conn, err := jsonrpc.Dial(addr, c.server.IntCACert, c.cert, c.key, false) // no need for debug mode on conn as we have it on server conn already
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

		conn.SetReadDeadline(10)
		c.conn = conn
		return c
	}

	return nil
}

// WriteRequest sends a request message through the connection.
func (c *ConnHelper) WriteRequest(method string, params interface{}) (reqID string) {
	id, err := c.conn.WriteRequest(method, params)
	if err != nil {
		c.testing.Fatal("Failed to write request to connection:", err)
	}
	return id
}

// WriteRequestArr sends a request message through the connection. Params object is variadic.
func (c *ConnHelper) WriteRequestArr(method string, params ...interface{}) (reqID string) {
	return c.WriteRequest(method, params)
}

// WriteNotification sends a notification message through the connection.
func (c *ConnHelper) WriteNotification(method string, params interface{}) {
	if err := c.conn.WriteNotification(method, params); err != nil {
		c.testing.Fatal("Failed to write notification to connection:", err)
	}
}

// WriteNotificationArr sends a notification message through the connection. Params object is variadic.
func (c *ConnHelper) WriteNotificationArr(method string, params ...interface{}) {
	c.WriteNotification(method, params)
}

// WriteResponse sends a response message through the connection.
func (c *ConnHelper) WriteResponse(id string, result interface{}, err *jsonrpc.ResError) {
	if err := c.conn.WriteResponse(id, result, err); err != nil {
		c.testing.Fatal("Failed to write response to connection:", err)
	}
}

// ReadReq reads a request object from a connection.
// If incoming message is not a request, a fatal error is logged.
// Optionally, you can pass in a data structure that the returned JSON-RPC request params data will be serialized into.
// Otherwise json.Unmarshal defaults apply.
func (c *ConnHelper) ReadReq(paramsData interface{}) *jsonrpc.Request {
	req, _, _, err := c.conn.ReadMsg(nil, paramsData)
	if err != nil {
		c.testing.Fatal("Failed to read request from connection:", err)
	}

	if req == nil {
		c.testing.Fatal("Read message was not a request message.")
	}

	return req
}

// ReadRes reads a response object from a connection.
// If incoming message is not a response, a fatal error is logged.
// Optionally, you can pass in a data structure that the returned JSON-RPC response result data will be serialized into.
// Otherwise json.Unmarshal defaults apply.
func (c *ConnHelper) ReadRes(resultData interface{}) *jsonrpc.Response {
	_, res, _, err := c.conn.ReadMsg(resultData, nil)
	if err != nil {
		c.testing.Fatal("Failed to read response from connection:", err)
	}

	if res == nil {
		c.testing.Fatal("Read message was not a response message.")
	}

	return res
}

// ReadMsg reads a JSON-RPC message from a connection. If possible, use more specific ReadReq/ReadRes/ReadNot methods instead.
// Optionally, you can pass in a data structure that the returned JSON-RPC response result data will be serialized into (same for request params).
// Otherwise json.Unmarshal defaults apply.
func (c *ConnHelper) ReadMsg(resultData interface{}, paramsData interface{}) (req *jsonrpc.Request, res *jsonrpc.Response, not *jsonrpc.Notification) {
	req, res, not, err := c.conn.ReadMsg(resultData, paramsData)
	if err != nil {
		c.testing.Fatal("Failed to read message from connection:", err)
	}

	return
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

// Close closes a connection.
func (c *ConnHelper) Close() {
	if err := c.conn.Close(); err != nil {
		c.testing.Fatal("Failed to close connection:", err)
	}
}
