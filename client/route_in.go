package client

// ------ Incoming Requests ---------- //

// InMsgHandler registers a handler to accept incoming messages from the server.
func (c *Client) InMsgHandler(handler func(m []Message) error) {
	c.inMsgHandler = handler
}
