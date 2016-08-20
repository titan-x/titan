package client

import (
	"fmt"

	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
	"github.com/titan-x/titan/models"
)

// ------ Incoming Requests ---------- //

// InMsgHandler registers a handler to accept incoming messages from the server.
func (c *Client) InMsgHandler(handler func(m []models.Message) error) {
	r := middleware.NewRouter()
	c.conn.Middleware(r)
	r.Request("msg.recv", func(ctx *neptulon.ReqCtx) error {
		var msg []models.Message
		if err := ctx.Params(&msg); err != nil {
			return fmt.Errorf("client: msg.recv: error reading request params: %v", err)
		}

		if err := handler(msg); err != nil {
			return err
		}

		ctx.Res = ACK
		return ctx.Next()
	})
}
