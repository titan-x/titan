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
		to      string
		message string
	}

	return func(ctx *jsonrpc.ReqCtx) {
		var r sendMsgReq
		ctx.Params(&r)
		q.AddRequest(r.to, "msg.recv", r.message, func(ctx *jsonrpc.ResCtx) {
			// todo: send 'delivered' message to sender about this message (or failed depending on output)
		})

		ctx.Res = "ACK"
	}
}

func initRecvMsgHandler() func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {

	}
}
