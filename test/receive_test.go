package test

import "testing"

func TestReceiveOfflineQueue(t *testing.T) {
	s := NewServerHelper(t).SeedDB()
	defer s.Stop()
	c1 := NewClientHelper(t).DefaultCert().Dial()
	defer c1.Close()

	_ = c1.WriteRequest("msg.send", nil)

	c2 := NewClientHelper(t).Cert(client2Cert, client2Key).Dial()
	defer c2.Close()

	// _ = c1.WriteRequest("msg.recv", nil)

	// t.Fatal("Failed to receive queued messages after coming online")
	// t.Fatal("Failed to send ACK for received message queue")
}

func TestReceiveEcho(t *testing.T) {
	// send message to user with ID: "client.127.0.0.1"
}
