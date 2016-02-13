package titan

import (
	"time"

	"github.com/neptulon/cmap"
	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
)

// Client is a Titan client.
type Client struct {
	ID           string     // Randomly generated unique client connection ID.
	Session      *cmap.CMap // Thread-safe data store for storing arbitrary data for this connection session.
	conn         *neptulon.Conn
	router       *middleware.Router
	inMsgHandler func(m []Message) error
}

// NewClient creates a new Client object.
func NewClient() (*Client, error) {
	c, err := neptulon.NewConn()
	if err != nil {
		return nil, err
	}
	r := middleware.NewRouter()
	c.Middleware(r)
	s := &Client{ID: c.ID, Session: c.Session, conn: c, router: r, inMsgHandler: func(m []Message) error { return nil }}
	r.Request("msg.recv", s.inMsgRoute)
	return s, nil
}

// SetDeadline set the read/write deadlines for the connection, in seconds.
func (c *Client) SetDeadline(seconds int) {
	c.conn.SetDeadline(seconds)
}

// Connect connectes to the server at given network address and starts receiving messages.
func (c *Client) Connect(addr string) error {
	return c.conn.Connect(addr)
}

// Close closes a client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) inMsgRoute(ctx *neptulon.ReqCtx) error {
	var msg []Message
	if err := ctx.Params(msg); err != nil {
		return err
	}

	if err := c.inMsgHandler(msg); err != nil {
		return err
	}

	return ctx.Next()
}

// ---- In/Out Request Objects ------ //

// Message is a chat message.
type Message struct {
	From    string    `json:"from,omitempty"`
	To      string    `json:"to,omitempty"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// ------ Incoming Requests ---------- //

// InMsgHandler registers a handler to accept incoming messages from the server.
func (c *Client) InMsgHandler(handler func(m []Message) error) {
	c.inMsgHandler = handler
}

// ------ Outgoing Requests ---------- //

// GetPendingMessages sends a request to server to receive any pending messages.
func (c *Client) GetPendingMessages(msgHandler func(m []Message) error) error {
	_, err := c.conn.SendRequest("msg.recv", nil, func(ctx *neptulon.ResCtx) error {
		var msg []Message
		if err := ctx.Result(msg); err != nil {
			return err
		}
		if err := msgHandler(msg); err != nil {
			return err
		}
		return nil
	})

	return err
}

// SendMessages sends a batch of messages to the server.
func (c *Client) SendMessages(m []Message) error {
	_, err := c.conn.SendRequest("msg.send", m, func(ctx *neptulon.ResCtx) error {
		return nil
	})

	return err
}

// Echo sends a message to server echo endpoint.
// This is meant to be used for testing connectivity.
func (c *Client) Echo(m interface{}, msgHandler func(msg *Message) error) error {
	_, err := c.conn.SendRequest("msg.echo", m, func(ctx *neptulon.ResCtx) error {
		var msg Message
		if err := ctx.Result(&msg); err != nil {
			return err
		}
		if err := msgHandler(&msg); err != nil {
			return err
		}
		return nil
	})

	return err
}
