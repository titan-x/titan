package titan

import (
	"sync"

	"github.com/neptulon/cmap"
	"github.com/neptulon/jsonrpc"
	"github.com/neptulon/neptulon"
)

// Client is a Titan API client.
type Client struct {
	client *jsonrpc.Client
}

// NewClient creates a new Client object.
// msgWG = (optional) sets the given *sync.WaitGroup reference to be used for counting active gorotuines that are used for handling incoming/outgoing messages.
// disconnHandler = (optional) registers a function to handle client disconnection events.
func NewClient(msgWG *sync.WaitGroup, disconnHandler func(client *neptulon.Client)) *Client {
	return &Client{client: jsonrpc.NewClient(msgWG, disconnHandler)}
}

// Connect connectes to the server at given network address and starts receiving messages.
func (c *Client) Connect(addr string, debug bool) error {
	return c.client.Connect(addr, debug)
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

// Close closes a client connection.
func (c *Client) Close() error {
	return c.client.Close()
}
