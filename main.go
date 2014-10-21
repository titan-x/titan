package main

import (
	"fmt"
	"log"

	"github.com/soygul/nbusy-server/gcm/ccs"
)

func main() {
	config := GetConfig()
	// config.App.Debug = true
	ccsConn, err := ccs.New(config.GCM.CCSEndpoint, config.GCM.SenderID, config.GCM.APIKey, config.App.Debug)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully logged in to GCM.")

	msgCh := make(chan map[string]interface{})
	errCh := make(chan error)

	go ccsConn.Listen(msgCh, errCh)

	ccsMessage := ccs.NewMessage(config.GCM.RegID)
	ccsMessage.SetData("hello", "world")
	ccsConn.Send(ccsMessage)

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
