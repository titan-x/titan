package ccs

import (
	"os"
	"testing"
)

// GCM environment variables
var host = os.Getenv("GCM_CCS_HOST")
var senderID = os.Getenv("GCM_SENDER_ID")
var apiKey = os.Getenv("GOOGLE_API_KEY")
var regID = os.Getenv("GCM_REG_ID") // optional registration ID from an Android device testing outgoing messages

func TestConnect(t *testing.T) {
	c, err := getConn(t)
	if err != nil {
		t.Fatal(err)
	}
	c.Close()
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
	} else if host == "" || senderID == "" || apiKey == "" {
		t.Skip("skipping integration test due to missing GCM environment variables.")
	}

	return Connect(host, senderID, apiKey, true)
}
