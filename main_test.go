package main

import (
	"testing"
	"fmt"
	"github.com/soygul/nbusy-server/gcm/ccs"
)

func TestMain(t *testing.T) {
	config := GetConfig()
	ccsConn, err := ccs.New(config.GCM.CCSEndpoint, config.GCM.SenderID, config.GCM.APIKey, config.App.Debug)
	if err != nil {
		t.Fatalf("Connection to CCS failed with error: %+v", err)
	}

	msgCh := make(chan map[string]interface{})
	errCh := make(chan error)

	go ccsConn.Listen(msgCh, errCh)

	ccsMessage := ccs.NewMessage("GCM_TEST_REG_ID")
	ccsMessage.SetData("hello", "world")
	ccsConn.Send(ccsMessage)

	for {
		select {
		case err := <-errCh:
			fmt.Println("err:", err)
			return;
		case msg := <-msgCh:
			fmt.Println("msg:", msg)
			return;
		}
	}
}
