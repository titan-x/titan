package devastator

import (
	"log"

	"github.com/nbusy/neptulon/jsonrpc"
)

func initPrivRoutes(r *jsonrpc.Router, q *Queue) {
	// r.Middleware(Logger()) - request-response logger (the pointer fields in request/response objects will have to change for this to work)
	r.Request("msg.echo", initEchoMsgHandler())
	r.Request("msg.send", initSendMsgHandler(q))
	r.Request("msg.recv", initRecvMsgHandler())
	// r.Middleware(NotFoundHandler()) - 404-like handler
}

func initEchoMsgHandler() func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {
		ctx.Params(&ctx.Res)
		if ctx.Res == nil {
			ctx.Res = ""
		}
	}
}

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
			// todo: send 'delivered' message to sender (as a request?) about this message (or failed, depending on output)
			var res string
			ctx.Result(&res)
			if res == "ACK" {
				// q.AddRequest(uid, "msg.delivered", ... // requeue if failed or handle resends automatically in the queue type, which is prefered)
			}
		})

		if err != nil {
			log.Fatal(err)
			return
		}

		ctx.Res = "ACK"
	}
}

func initRecvMsgHandler() func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {
		ctx.Res = "ACK"
	}
}
