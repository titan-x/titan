package main

import (
	"fmt"
	"log"
	"os"

	"github.com/soygul/nbusy-server/gcm/ccs"
)

const (
	GO_ENV = "development"
	GCM_CCS_ENDPOINT         = "gcm.googleapis.com:5235"
	GCM_CCS_STAGING_ENDPOINT = "gcm-staging.googleapis.com:5236"
	GCM_SENDER_ID            = ""
	GOOGLE_API_KEY           = ""
	GCM_TEST_REG_ID          = ""
)

// GCM type describes the Google Cloud Messaging parameters as described here: https://developer.android.com/google/gcm/gs.html
type GCM struct {
	CCSEndpoint string
	SenderID    string
	APIKey      string
}

func main() {
	gcm := GCM{CCSEndpoint: os.Getenv("GCM_CCS_ENDPOINT"), SenderID: os.Getenv("GCM_SENDER_ID"), APIKey: os.Getenv("GOOGLE_API_KEY")}



	ccsClient, err := ccs.New(GCM_SENDER_ID, GOOGLE_API_KEY, false)
	if err != nil {
		log.Fatal(err)
	}

	msgCh := make(chan map[string]interface{})
	errCh := make(chan error)

	go ccsClient.Recv(msgCh, errCh)

	ccsMessage := ccs.NewMessage("GCM_TEST_REG_ID")
	ccsMessage.SetData("hello", "world")
	ccsMessage.CollapseKey = ""
	ccsMessage.TimeToLive = 0
	ccsMessage.DelayWhileIdle = true
	ccsClient.Send(ccsMessage)

	fmt.Print("NBusy messege server started.")

	for {
		select {
		case err := <-errCh:
			fmt.Println("err:", err)
		case msg := <-msgCh:
			fmt.Println("msg:", msg)
		}
	}
}
