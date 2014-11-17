package main

import (
	"log"

	"github.com/soygul/nbusy-server/gcm/ccs"
)

func main() {
	conf := GetConfig()
	c, err := ccs.Connect(conf.GCM.CCSEndpoint, conf.GCM.SenderID, conf.GCM.APIKey, conf.App.Debug)

	if conf.App.Debug {
		if err == nil {
			log.Printf("New CCS connection established with parameters: %+v\n", c)
		} else {
			log.Fatalf("New CCS connection failed to establish with parameters: %+v\n", c)
		}
	}

	log.Println("NBusy messege server started.")

	for {
		m, err := c.Receive()
		if err != nil {
			log.Printf("Incoming CCS error: %v\n", err)
		}
		go readHandler(m)
	}
}

func readHandler(m *ccs.InMsg) {
	log.Printf("Incoming CCS message: %v\n", m)
}
