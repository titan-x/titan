package client

import (
	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
)

const (
	// ACK is the short acknowledgement response for a request.
	ACK = "ACK"

	// NACK is the short rejection response for a request.
	NACK = "NACK"
)

// Client is a Titan client.
type Client struct {
	ID      string     // Randomly generated unique client connection ID.
	Session *cmap.CMap // Thread-safe data store for storing arbitrary data for this connection session.
	conn    *neptulon.Conn
	router  *middleware.Router
}

// NewClient creates a new Client object.
func NewClient() (*Client, error) {
	c, err := neptulon.NewConn()
	if err != nil {
		return nil, err
	}

	r := middleware.NewRouter()
	c.Middleware(r)

	return &Client{
		ID:      c.ID,
		Session: c.Session,
		conn:    c,
		router:  r,
	}, nil
}

// SetDeadline set the read/write deadlines for the connection, in seconds.
func (c *Client) SetDeadline(seconds int) {
	c.conn.SetDeadline(seconds)
}

// Middleware registers middleware to handle incoming request messages.
func (c *Client) Middleware(middleware ...neptulon.Middleware) {
	c.conn.Middleware(middleware...)
}

// MiddlewareFunc registers middleware function to handle incoming request messages.
func (c *Client) MiddlewareFunc(middleware ...func(ctx *neptulon.ReqCtx) error) {
	c.conn.MiddlewareFunc(middleware...)
}

// DisconnHandler registers a function to handle disconnection event.
func (c *Client) DisconnHandler(handler func(c *Client)) {
	c.conn.DisconnHandler(func(nepc *neptulon.Conn) {
		handler(c)
	})
}

// Connect connectes to the server at given network address and starts receiving messages.
func (c *Client) Connect(addr string) error {
	return c.conn.Connect(addr)
}

// Close closes a client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
