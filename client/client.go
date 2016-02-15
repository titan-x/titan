package client

import (
	"fmt"
	"sync"

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
	jwtToken     string
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

// UseJWT enables JWT authentication.
func (c *Client) UseJWT(token string) {
	c.jwtToken = token
}

// DisconnHandler registers a function to handle disconnection event.
func (c *Client) DisconnHandler(handler func(c *Client)) {
	c.conn.DisconnHandler(func(nepc *neptulon.Conn) {
		handler(c)
	})
}

// Connect connectes to the server at given network address and starts receiving messages.
func (c *Client) Connect(addr string) error {
	if err := c.conn.Connect(addr); err != nil {
		return err
	}

	if c.jwtToken != "" {
		var wg sync.WaitGroup
		wg.Add(1)
		if err := c.jwtAuth(c.jwtToken, func(ack string) error {
			defer wg.Done()
			if ack != "ACK" {
				return fmt.Errorf("server did not ACK our auth.jwt request: %v", ack)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("authentication failed: %v", err)
		}
		wg.Wait()
	}

	return nil
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

	ctx.Res = "ACK"
	return ctx.Next()
}
