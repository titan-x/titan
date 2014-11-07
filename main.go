package main

import (
	"fmt"
	"log"

	"github.com/soygul/nbusy-server/gcm/ccs"
)

func main() {
	config := GetConfig()
	conn, err := ccs.Connect(config.GCM.CCSEndpoint, config.GCM.SenderID, config.GCM.APIKey, config.App.Debug)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully logged in to GCM.")

	go conn.Listen()

	ccsMessage := ccs.NewMessage(config.GCM.RegID)
	ccsMessage.SetData("hello", "world")
	conn.Send(ccsMessage)

	fmt.Println("NBusy messege server started.")

	for {
		select {
		case err := <-conn.ErrorChan:
			fmt.Println("err:", err)
		case msg := <-conn.MessageChan:
			fmt.Println("msg:", msg)
		}
	}
}
