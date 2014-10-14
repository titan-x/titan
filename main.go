package main

import (
	"fmt"
	"log"
	"github.com/soygul/nbusy-server/gcm/ccs"
)

func main() {
	config := GetConfig()
	config.App.Debug = true;
	ccsConnection, err := ccs.New(config.GCM.CCSEndpoint, config.GCM.SenderID, config.GCM.APIKey, config.App.Debug)
	fmt.Println(ccsConnection)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully logged in to GCM.")

	msgCh := make(chan map[string]interface{})
	errCh := make(chan error)

	go ccsConnection.Listen(msgCh, errCh)

//	ccsMessage := ccs.NewMessage("GCM_TEST_REG_ID")
//	ccsMessage.SetData("hello", "world")
//	ccsMessage.CollapseKey = ""
//	ccsMessage.TimeToLive = 0
//	ccsMessage.DelayWhileIdle = true
//	ccsClient.Send(ccsMessage)

	fmt.Println("NBusy messege server started.")

	for {
		select {
		case err := <-errCh:
			fmt.Println("err:", err)
		case msg := <-msgCh:
			fmt.Println("msg:", msg)
		}
	}
}
