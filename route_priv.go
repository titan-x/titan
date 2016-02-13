package titan

import (
	"fmt"

	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
	"github.com/titan-x/titan/client"
)

func initPrivRoutes(r *middleware.Router, q *Queue) {
	r.Request("msg.echo", middleware.Echo)
	r.Request("msg.send", initSendMsgHandler(q))
	r.Request("msg.recv", initRecvMsgHandler(q))
}

// Allows clients to send messages to each other, online or offline.
func initSendMsgHandler(q *Queue) func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		var s client.Message
		ctx.Params(&s)

		uid := ctx.Conn.Session.Get("userid").(string)
		r := client.Message{From: uid, Message: s.Message}
		err := q.AddRequest(s.To, "msg.recv", r, func(ctx *neptulon.ResCtx) error {
			var res string
			ctx.Result(&res)
			if res == "ACK" {
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

		ctx.Res = "ACK"
		return ctx.Next()
	}
}

// Used only for a client to announce its presence.
// If there are any messages meant for this user, they are started to be sent with this call (via the cert-auth middleware).
func initRecvMsgHandler(q *Queue) func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		q.SetConn(ctx.Conn.Session.Get("userid").(string), ctx.Conn.ID)
		ctx.Res = "ACK" // todo: this could rather send the remaining queue size for the client
		return ctx.Next()
	}
}
