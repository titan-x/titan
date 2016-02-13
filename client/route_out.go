package client

import "github.com/neptulon/neptulon"

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
