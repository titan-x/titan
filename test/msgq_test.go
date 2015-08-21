package test

import "testing"

func TestReceiveQueue(t *testing.T) {
	s := NewServerHelper(t).SeedDB()
	defer s.Stop()
	c1 := NewClientHelper(t).DefaultCert().Dial()
	defer c1.Close()

	_ = c1.WriteRequest("msg.send", nil)

	c2 := NewClientHelper(t).Cert(client2Cert, client2Key).Dial()
	defer c2.Close()

	_ = c1.WriteRequest("msg.recv", nil)
}
