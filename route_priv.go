package devastator

import "github.com/nbusy/neptulon/jsonrpc"

func initPrivRoutes(r *jsonrpc.Router) {

}

func initMsgSendHandler(q *Queue) func(ctx *jsonrpc.ReqCtx) {
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
