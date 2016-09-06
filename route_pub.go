package titan

import (
	"github.com/neptulon/neptulon"
	"github.com/neptulon/neptulon/middleware"
	"github.com/titan-x/titan/data"
)

// we need *data.DB (pointer to interface) so that the closure below won't capture the actual value that pointer points to
// so we can swap databases whenever we want using Server.SetDB(...)
func initPubRoutes(r *middleware.Router, db *data.DB, pass string) {
	r.Request("auth.google", initGoogleAuthHandler(db, pass))
}

func initGoogleAuthHandler(db *data.DB, pass string) func(ctx *neptulon.ReqCtx) error {
	return func(ctx *neptulon.ReqCtx) error {
		if err := googleAuth(ctx, *db, pass); err != nil {
			return err
		}
		return nil
	}
}
