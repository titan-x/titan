package test

import (
	"testing"

	"github.com/titan-x/titan/data"
)

func TestEchoBot(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	ch := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect().JWTAuthSync()
	defer ch.CloseWait()

	ch.EchoSync("Ola!") // todo: ch.EchoBot("Ola!")
}
