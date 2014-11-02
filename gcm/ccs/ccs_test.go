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

func getConn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short testing mode.")
	} else if senderID == "" || regID == "" || ccsEndpoint == "" || apiKey == "" {
		t.Skip("skipping integration test due to missing GCM configuration environment variables.")
	}

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
