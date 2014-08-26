package main

import (
	"fmt"
	"log"
	"os"

	"github.com/soygul/nbusy-server/gcm/ccs"
)

func main() {
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
