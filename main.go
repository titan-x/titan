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

	ccsMessage := ccs.NewMessage(config.GCM.RegID)
	ccsMessage.SetData("hello", "world")
	conn.Send(ccsMessage)

	fmt.Println("NBusy messege server started.")

	for {
		msg, err := conn.Read()
		if err != nil {
			log.Printf("Incoming CCS error: %v\n", err)
		}
		go readHandler(msg)
	}
}

func readHandler(msg map[string]interface{}) {
	log.Printf("Incoming CCS message: %v\n", msg)
}
