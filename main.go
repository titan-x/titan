package main

import (
	"log"

	"github.com/soygul/gcm-ccs"
)

func main() {
	c, err := ccs.Connect(Conf.GCM.CCSHost, Conf.GCM.SenderID, Conf.GCM.APIKey, Conf.App.Debug)
	if err != nil {
		log.Fatalf("NBusy messege server failed to start.")
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
