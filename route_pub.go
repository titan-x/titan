package titan

import (
	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
)

func initPubRoutes(r *middleware.Router, db DB) {
	r.Request("auth.google", initGoogleAuthHandler(db))
	// pubRoute.NotFound(...)
	// todo: if the first incoming message in public route is not one of close/google.auth,
	// close the connection right away (and maybe wait for client to return ACK then close?)
}

func initGoogleAuthHandler(db DB) func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		if err := googleAuth(ctx, db); err != nil {
			return err
		}
		return ctx.Next()
	}
}
