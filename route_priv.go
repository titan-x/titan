package titan

import (
	"fmt"
	"strings"

	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
	"github.com/titan-x/titan/client"
	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/models"
)

// We need *data.Queue (pointer to interface) so that the closure below won't capture the actual value that pointer points to
// so we can swap queues whenever we want using Server.SetQueue(...)
func initPrivRoutes(r *middleware.Router, q *data.Queue) {
	r.Request("auth.jwt", initJWTAuthHandler())
	r.Request("echo", middleware.Echo)
	r.Request("msg.send", initSendMsgHandler(q))
}

// Used for a client to authenticate and announce its presence.
// If there are any messages meant for this user, they are started to be sent after this call.
func initJWTAuthHandler() func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		// todo: this could rather send the remaining queue size for the client so client can disconnect if there is nothing else to do
		ctx.Res = client.ACK
		return ctx.Next()
	}
}

// Allows clients to send messages to each other, online or offline.
func initSendMsgHandler(q *data.Queue) func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		var sMsgs []models.Message
		if err := ctx.Params(&sMsgs); err != nil {
			return err
		}

		uid := ctx.Conn.Session.Get("userid").(string)

		for _, sMsg := range sMsgs {
			from := uid
			to := strings.ToLower(sMsg.To)

			// handle messages to bots
			if to == "echo" {
				from = "echo"
				to = uid
			}

			// submit the messages to send queue
			err := (*q).AddRequest(to, "msg.recv", []models.Message{models.Message{From: from, Message: sMsg.Message}}, func(ctx *neptulon.ResCtx) error {
				var res string
				ctx.Result(&res)
				if res == client.ACK {
					// todo: send 'delivered' message to sender (as a request?) about this message (or failed, depending on output)
					// todo: q.AddRequest(uid, "msg.delivered", ... // requeue if failed or handle resends automatically in the queue type, which is prefered)
				} else {
					// todo: auto retry or "msg.failed" ?
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("route: msg.recv: failed to add request to queue with error: %v", err)
			}
		}

		ctx.Res = client.ACK
		return ctx.Next()
	}
}
