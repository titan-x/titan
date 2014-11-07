package ccs

import (
	"os"
	"testing"
)

// GCM environment variables
var senderID = os.Getenv("GCM_SENDER_ID")
var regID = os.Getenv("GCM_REG_ID")
var ccsEndpoint = os.Getenv("GCM_CCS_ENDPOINT")
var apiKey = os.Getenv("GOOGLE_API_KEY")

func TestConnect(t *testing.T) {
}

func TestDisconnect(t *testing.T) {
}

func TestGCMMessages(t *testing.T) {
	// see if we can handle all known GCM message types properly
}

func TestMessageFields(t *testing.T) {
	// see if our message structure's fields match the incoming message fields exactly
}

func TestReceive(t *testing.T) {
}

func TestSend(t *testing.T) {
}

func getConn(t *testing.T) (Conn, error) {
	if testing.Short() {
		t.Skip("skipping integration test in short testing mode.")
	} else if senderID == "" || regID == "" || ccsEndpoint == "" || apiKey == "" {
		t.Skip("skipping integration test due to missing GCM configuration environment variables.")
	}

	return Conn{}, nil
}
