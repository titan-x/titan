package main

import (
	"fmt"
	"log"
)

const (
	GCM_CLIENT_ID = ""
	GCM_API_KEY   = ""

	TEST_REG_ID = ""
)

func main() {
	ccsClient, err := ccs.New(GCM_CLIENT_ID, GCM_API_KEY, false)
	if err != nil {
		log.Fatal(err)
	}

	msgCh := make(chan map[string]interface{})
	errCh := make(chan error)

	go ccsClient.Recv(msgCh, errCh)

	ccsMessage := ccs.NewMessage(TEST_REG_ID)
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
