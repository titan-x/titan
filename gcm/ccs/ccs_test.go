package ccs

import "testing"

func TestConnect(t *testing.T) {
	t.SkipNow()
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
