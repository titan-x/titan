package main

import (
	"log"

	"github.com/soygul/nbusy-server/gcm/ccs"
)

func main() {
	c, err := ccs.Connect(Conf.GCM.CCSHost, Conf.GCM.SenderID, Conf.GCM.APIKey, Conf.App.Debug)

	if Conf.App.Debug {
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
