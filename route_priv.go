package devastator

import "github.com/nbusy/neptulon/jsonrpc"

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

		r := recvMsgReq{From: ctx.Conn.Data.Get("userid").(string), Message: s.Message}
		err := q.AddRequest(s.To, "msg.recv", r, func(ctx *jsonrpc.ResCtx) {
			// todo: send 'delivered' message to sender (as a request?) about this message (or failed depending on output)
		})

		if err != nil {
			// todo: return error or "NACK" ?
			return
		}

		ctx.Res = "ACK"
	}
}

func initRecvMsgHandler() func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {

	}
}
