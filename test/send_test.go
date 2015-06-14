package test

import "testing"

func TestSendEcho(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t, true)
	defer c.Close()

	id := c.WriteRequest("echo", map[string]string{"echo": "echo"})
	m := c.ReadMsg()

	res := m.Result.(map[string]interface{})
	if m.ID != id || res["echo"] != "echo" {
		t.Fatal("Failed to receive echo message in proper format:", m)
	}

	// t.Fatal("Failed to send a message to the echo user")
	// t.Fatal("Failed to send batch message to the echo user")
	// t.Fatal("Failed to send large message to the echo user")
	// t.Fatal("Did not receive ACK for a message sent")
	// t.Fatal("Failed to receive a response from echo user")
	// t.Fatal("Could not send an ACK for an incoming message")
}

func TestPing(t *testing.T) {
	// t.Fatal("Pong/ACK was not sent for ping")
}
