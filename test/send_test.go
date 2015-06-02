package test

import "testing"

func TestSendEcho(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)

	closeClientConn(t, c)
	stopServer(t, s)

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
