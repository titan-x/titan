package devastator

import (
	"log"

	"github.com/nbusy/neptulon/jsonrpc"
)

func initPrivRoutes(r *jsonrpc.Router, q *Queue) {
	// r.Middleware(Logger()) - request-response logger (the pointer fields in request/response objects will have to change for this to work)
	r.Request("msg.echo", initEchoMsgHandler())
	r.Request("msg.send", initSendMsgHandler(q))
	r.Request("msg.recv", initRecvMsgHandler(q))
	// r.Middleware(NotFoundHandler()) - 404-like handler
}

// Echoes message sent by the client back to the client.
func initEchoMsgHandler() func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {
		ctx.Params(&ctx.Res)
		if ctx.Res == nil {
			ctx.Res = ""
		}
	}
}

// Allows clients to send messages to each other, online or offline.
func initSendMsgHandler(q *Queue) func(ctx *jsonrpc.ReqCtx) {
	type sendMsgReq struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}

	type recvMsgReq struct {
		From    string `json:"from"`
		Message string `json:"message"`
	}

	return func(ctx *jsonrpc.ReqCtx) {
		var s sendMsgReq
		ctx.Params(&s)

		uid := ctx.Conn.Data.Get("userid").(string)
		r := recvMsgReq{From: uid, Message: s.Message}
		err := q.AddRequest(s.To, "msg.recv", r, func(ctx *jsonrpc.ResCtx) {
			var res string
			ctx.Result(&res)
			if res == "ACK" {
				// todo: send 'delivered' message to sender (as a request?) about this message (or failed, depending on output)
				// todo: q.AddRequest(uid, "msg.delivered", ... // requeue if failed or handle resends automatically in the queue type, which is prefered)
			} else {
				// todo: auto retry or "msg.failed" ?
			}
		})

		if err != nil {
			log.Fatal("Failed to add request to queue with error:", err)
			return
		}

		ctx.Res = "ACK"
	}
}

// Used only for a client to announce its presence.
// If there are any messages meant for this user, they are started to be sent with this call (via the cert-auth middleware).
func initRecvMsgHandler(q *Queue) func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {
		q.SetConn(ctx.Conn.Data.Get("userid").(string), ctx.Conn.ID)
		ctx.Res = "ACK"
	}
}
