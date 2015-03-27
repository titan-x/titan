package main

import "testing"

func TestGoogleAuth(t *testing.T) {
	// t.Fatal("Google+ sign-in failed with valid credentials.")
	// t.Fatal("Google+ sign-in passed with invalid credentials.")
}

func TestClientCertAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid client certificate.")
	// t.Fatal("Authenticated with invalid client certificate.")
}

func TestStop(t *testing.T) {
	// todo: this should be a listener/queue test if we don't use any goroutines in the Server struct methods.
	// t.Fatal("Failed to stop the server gracefully: not all the goroutines were terminated properly.")
	// t.Fatal("Failed to stop the server gracefully: server did not wait for ongoing read/write operations.")
}
