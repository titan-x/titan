package middleware

import (
	"log"

	"github.com/neptulon/neptulon"
)

// Logger is an incoming/outgoing message logger.
func Logger(ctx *neptulon.ReqCtx) error {
	var v interface{}
	if err := ctx.Params(&v); err != nil {
		return err
	}

	err := ctx.Next()
	var res = ctx.Res
	if res == nil {
		res = ctx.Err
	}
	log.Printf("mw: logger: %v: %v, in: \"%v\", out: \"%#v\"", ctx.ID, ctx.Method, v, res)
	return err
}
