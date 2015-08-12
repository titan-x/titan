package test

import "testing"

func TestReceiveQueue(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).DefaultCert().Dial()
	defer c.Close()
}
