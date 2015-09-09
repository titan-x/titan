package devastator

import "github.com/nbusy/neptulon/jsonrpc"

func initPrivRoutes(r *jsonrpc.Router, q *Queue) {
	r.Request("msg.echo", initEchoMsgHandler())
	r.Request("msg.send", initSendMsgHandler(q))
	r.Request("msg.recv", initRecvMsgHandler())
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
			// todo: auto handle response (which is probably just ACK so this might be automated as a part of Queue or with something like MsgHandler)
			// or q.AddRequest(..).ACK();
		})
	}
}

func initRecvMsgHandler() func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {

	}
}
