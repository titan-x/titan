package titan

import (
	"fmt"

	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
)

func initPrivRoutes(r *middleware.Router, q *Queue) {
	// r.Middleware(Logger()) - request-response logger (the pointer fields in request/response objects will have to change for this to work)
	r.Request("msg.echo", initEchoMsgHandler())
	r.Request("msg.send", initSendMsgHandler(q))
	r.Request("msg.recv", initRecvMsgHandler(q))
	// r.Middleware(NotFoundHandler()) - 404-like handler
}

// Echoes message sent by the client back to the client.
func initEchoMsgHandler() func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		ctx.Params(&ctx.Res)
		if ctx.Res == nil {
			ctx.Res = ""
		}
		return ctx.Next()
	}
}

// Allows clients to send messages to each other, online or offline.
func initSendMsgHandler(q *Queue) func(ctx *neptulon.ReqCtx) error {
	type sendMsgReq struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}

	type recvMsgReq struct {
		From    string `json:"from"`
		Message string `json:"message"`
	}

	return func(ctx *neptulon.ReqCtx) error {
		var s sendMsgReq
		ctx.Params(&s)

		uid := ctx.Conn.Session.Get("userid").(string)
		r := recvMsgReq{From: uid, Message: s.Message}
		err := q.AddRequest(s.To, "msg.recv", r, func(ctx *neptulon.ResCtx) error {
			var res string
			ctx.Result(&res)
			if res == "ACK" {
				// todo: send 'delivered' message to sender (as a request?) about this message (or failed, depending on output)
				// todo: q.AddRequest(uid, "msg.delivered", ... // requeue if failed or handle resends automatically in the queue type, which is prefered)
			} else {
				// todo: auto retry or "msg.failed" ?
			}
			return ctx.Next()
		})

		if err != nil {
			return fmt.Errorf("Failed to add request to queue with error: %v", err)
		}

		ctx.Res = "ACK"
		return ctx.Next()
	}
}

// Used only for a client to announce its presence.
// If there are any messages meant for this user, they are started to be sent with this call (via the cert-auth middleware).
func initRecvMsgHandler(q *Queue) func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		q.SetConn(ctx.Client.Session().Get("userid").(string), ctx.Client.ConnID())
		ctx.Res = "ACK" // todo: this could rather send the remaining queue size for the client
		return ctx.Next()
	}
}
