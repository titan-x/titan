package main

import (
	"log"

	"github.com/soygul/nbusy-server/gcm/ccs"
)

func main() {
	conf := GetConfig()
	c, err := ccs.Connect(conf.GCM.CCSEndpoint, conf.GCM.SenderID, conf.GCM.APIKey, conf.App.Debug)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully logged in to GCM.")

	m := ccs.NewMessage(conf.GCM.RegID)
	m.SetData("hello", "world")
	c.Send(m)

	log.Println("NBusy messege server started.")

	for {
		m, err := c.Read()
		if err != nil {
			log.Printf("Incoming CCS error: %v\n", err)
		}
		go readHandler(m)
	}
}

func readHandler(m map[string]string) {
	log.Printf("Incoming CCS message: %v\n", m)
}
