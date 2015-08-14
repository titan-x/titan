package test

import (
	"testing"

	"github.com/nbusy/devastator"
)

func TestReceiveQueue(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).DefaultCert().Dial()
	defer c.Close()

	db := s.GetDB()
	db.SaveUser(&devastator.User{ID: 1, Cert: certChain.ClientCert})
	db.SaveUser(&devastator.User{ID: 2})
}
