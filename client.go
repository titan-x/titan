package titan

import (
	"sync"
	"time"

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
	c := jsonrpc.NewClient(msgWG, disconnHandler)
	c.HandleRequest("msg.recv", nil) // todo: use common handler. shall the handlers be obligatory during construction? or we might expect an interface implementation?
	return &Client{client: c}
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

// ------ incoming message handlers ---------- //

// ------ outgoing message senders ---------- //

// RecvMsg is an incoming chat message.
type RecvMsg struct {
	From    string    `json:"from"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// GetPendingMessages sends a request to server to receive any pending messages.
func (c *Client) GetPendingMessages(msgHandler func(m []RecvMsg) error) error {
	_, err := c.client.SendRequest("msg.recv", nil, func(ctx *jsonrpc.ResCtx) error {
		var msg []RecvMsg
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
