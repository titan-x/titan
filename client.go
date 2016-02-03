package titan

import (
	"time"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon"
)

// Client is a Titan client.
type Client struct {
	conn *neptulon.Conn
}

// NewClient creates a new Client object.
func NewClient() *Client {
	return &Client{conn: neptulon.NewConn()}
}

// ConnID is a randomly generated unique client connection ID.
func (c *Client) ConnID() string {
	return c.client.ConnID()
}

// Session is a thread-safe data store for storing arbitrary data for this connection session.
func (c *Client) Session() *cmap.CMap {
	return c.client.Session()
}

// SetDeadline set the read/write deadlines for the connection, in seconds.
func (c *Client) SetDeadline(seconds int) {
	c.client.SetDeadline(seconds)
}

// UseTLS enables Transport Layer Security for the connection.
// ca = Optional CA certificate to be used for verifying the server certificate. Useful for using self-signed server certificates.
// clientCert, clientCertKey = Optional certificate/privat key pair for TLS client certificate authentication.
// All certificates/private keys are in PEM encoded X.509 format.
func (c *Client) UseTLS(ca, clientCert, clientCertKey []byte) {
	c.client.UseTLS(ca, clientCert, clientCertKey)
}

// Connect connectes to the server at given network address and starts receiving messages.
func (c *Client) Connect(addr string, debug bool) error {
	return c.client.Connect(addr, debug)
}

// Close closes a client connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// ---- In/Out Request Objects ------ //

// Message is a chat message.
type Message struct {
	From    string    `json:"from"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// ------ Incoming Requests ---------- //

// HandleIncomingMessages registers a handler to accept incoming messages from the server.
func (c *Client) HandleIncomingMessages(msgHandler func(m []Message) error) {
	c.client.HandleRequest("msg.recv", func(ctx *neptulon.ReqCtx) error {
		var msg []Message
		if err := ctx.Params(msg); err != nil {
			return err
		}

		if err := msgHandler(msg); err != nil {
			return err
		}

		return ctx.Next()
	})
}

// ------ Outgoing Requests ---------- //

// GetPendingMessages sends a request to server to receive any pending messages.
func (c *Client) GetPendingMessages(msgHandler func(m []Message) error) error {
	_, err := c.client.SendRequest("msg.recv", nil, func(ctx *neptulon.ResCtx) error {
		var msg []Message
		if err := ctx.Result(msg); err != nil {
			return err
		}
		if err := msgHandler(msg); err != nil {
			return err
		}
		return ctx.Next()
	})

	return err
}

// SendMessages sends a batch of messages to the server.
func (c *Client) SendMessages(m []Message) error {
	_, err := c.client.SendRequest("msg.send", m, func(ctx *neptulon.ResCtx) error {
		return ctx.Next()
	})

	return err
}

// Echo sends a message to server echo endpoint.
// This is meant to be used for testing connectivity.
func (c *Client) Echo(m interface{}, msgHandler func(msg *Message) error) error {
	_, err := c.client.SendRequest("msg.echo", m, func(ctx *neptulon.ResCtx) error {
		var msg Message
		if err := ctx.Result(&msg); err != nil {
			return err
		}
		if err := msgHandler(&msg); err != nil {
			return err
		}
		return ctx.Next()
	})

	return err
}
