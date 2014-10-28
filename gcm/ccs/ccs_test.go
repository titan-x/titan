package ccs

import (
	"os"
	"testing"
)

const (
	// GCM environment variables
	gcmSenderID    = "GCM_SENDER_ID"
	gcmRegID       = "GCM_REG_ID"
	gcmCcsEndpoint = "GCM_CCS_ENDPOINT"
	googleAPIKey   = "GOOGLE_API_KEY"
)

var senderID = os.Getenv(gcmSenderID)
var regID = os.Getenv(gcmRegID)
var ccsEndpoint = os.Getenv(gcmSenderID)
var apiKey = os.Getenv(googleAPIKey)

func TestConnect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
}

// recv message
// send message
// gcm message types
// message data fields match
// documentation descriptions match

func initServer(t *testing.T) {
	// config := GetConfig()
	// ccsConn, err := New(config.GCM.CCSEndpoint, config.GCM.SenderID, config.GCM.APIKey, config.App.Debug)
	// if err != nil {
	// 	t.Fatalf("Connection to CCS failed with error: %+v", err)
	// }
	// t.Log("CCS connection established.")
	//
	// msgCh := make(chan map[string]interface{})
	// errCh := make(chan error)
	//
	// go ccsConn.Listen(msgCh, errCh)
	//
	// ccsMessage := ccs.NewMessage(config.GCM.RegID)
	// ccsMessage.SetData("hello", "world")
	// ccsConn.Send(ccsMessage)
	//
	// for {
	// 	select {
	// 	case err := <-errCh:
	// 		fmt.Println("err:", err)
	// 		return
	// 	case msg := <-msgCh:
	// 		fmt.Println("msg:", msg)
	// 		return
	// 	}
	// }
}
