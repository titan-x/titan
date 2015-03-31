package main

import "testing"

func TestGoogleAuth(t *testing.T) {
	// t.Fatal("Google+ first sign-in (registration) failed with valid credentials")
	// t.Fatal("Google+ second sign-in (regular) failed with valid credentials")
	// t.Fatal("Google+ sign-in passed with invalid credentials")
	// t.Fatal("Authentication was not ACKed")
}

func TestClientCertAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid client certificate")
	// t.Fatal("Authenticated with invalid/expired client certificate")
	// t.Fatal("Authentication was not ACKed")
}

func TestTokenAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid token")
	// t.Fatal("Authenticated with invalid/expired token")
	// t.Fatal("Authentication was not ACKed")
}

func TestReceiveOfflineQueue(t *testing.T) {
	// t.Fatal("Failed to receive queued messages after coming online")
	// t.Fatal("Failed to send ACK for received message queue")
}

func TestSendEcho(t *testing.T) {
	// t.Fatal("Failed to send a message to the echo user")
	// t.Fatal("Failed to send batch message to the echo user")
	// t.Fatal("Failed to send large message to the echo user")
	// t.Fatal("Did not receive ACK for a message sent")
	// t.Fatal("Failed to receive a response from echo user")
	// t.Fatal("Could not send an ACK for an incoming message")
}

func TestStop(t *testing.T) {
	// todo: this should be a listener/queue test if we don't use any goroutines in the Server struct methods
	// t.Fatal("Failed to stop the server gracefully: not all the goroutines were terminated properly")
	// t.Fatal("Failed to stop the server gracefully: server did not wait for ongoing read/write operations")
}

// is this listener test? or do we handle these errors in the server?
func TestTimeout(t *testing.T) {
	// t.Fatal("Send timout did not occur")
	// t.Fatal("Wait timeout did not occur")
}

// same question here..
func TestDisconnect(t *testing.T) {
	// t.Fatal("Client method.close request was not handled properly")
	// t.Fatal("Client disconnect was not handled gracefully")
	// t.Fatal("Server method.close request was not handled properly")
	// t.Fatal("Server disconnect was not handled gracefully")
}

func TestPing(t *testing.T) {
	// t.Fatal("Pong/ACK was not sent for ping")
}
