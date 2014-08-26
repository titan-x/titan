package main

import (
	"fmt"
	"log"
	"os"
	"github.com/soygul/nbusy-server/gcm/ccs"
)

const (
	GCM_CCS_ENDPOINT = "gcm.googleapis.com:5235"
	GCM_CCS_STAGING_ENDPOINT = "gcm-staging.googleapis.com:5236"
	GCM_SENDER_ID = ""
	GOOGLE_API_KEY = ""
	GCM_TEST_REG_ID = ""
)

func main() {
	os.Getenv("GCM_CLIENT_ID")
	ccsClient, err := ccs.New(GCM_SENDER_ID, GOOGLE_API_KEY, false)
	if err != nil {
		log.Fatal(err)
	}

	msgCh := make(chan map[string]interface{})
	errCh := make(chan error)

	go ccsClient.Recv(msgCh, errCh)

	ccsMessage := ccs.NewMessage(GCM_TEST_REG_ID)
	ccsMessage.SetData("hello", "world")
	ccsMessage.CollapseKey = ""
	ccsMessage.TimeToLive = 0
	ccsMessage.DelayWhileIdle = true
	ccsClient.Send(ccsMessage)

	for {
		select {
		case err := <-errCh:
			fmt.Println("err:", err)
		case msg := <-msgCh:
			fmt.Println("msg:", msg)
		}
	}

	fmt.Print("NBusy messege server started.")
}
