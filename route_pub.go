package devastator

import "github.com/neptulon/jsonrpc"

func initPubRoutes(r *jsonrpc.Router, db DB, certMgr *CertMgr) {
	r.Request("auth.google", initGoogleAuthHandler(db, certMgr))
	r.Notification("conn.close", initCloseConnHandler())
	// pubRoute.NotFound(...)
	// todo: if the first incoming message in public route is not one of close/google.auth,
	// close the connection right away (and maybe wait for client to return ACK then close?)
}

func initGoogleAuthHandler(db DB, certMgr *CertMgr) func(ctx *jsonrpc.ReqCtx) {
	return func(ctx *jsonrpc.ReqCtx) {
		googleAuth(ctx, db, certMgr)
	}
}

func initCloseConnHandler() func(ctx *jsonrpc.NotCtx) {
	return func(ctx *jsonrpc.NotCtx) {
		ctx.Done = true
		ctx.Conn.Close()
	}
}
