package main

import (
	"log"

	"github.com/nbusy/gcm/ccs"
)

func main() {
	c, err := ccs.Connect(Conf.GCM.CCSHost, Conf.GCM.SenderID, Conf.GCM.getAPIKey(), Conf.App.Debug)
	if err != nil {
		log.Fatalln("Failed to connect to GCM CCS with error:", err)
	}

	log.Println("NBusy message server started.")

	for {
		m, err := c.Receive()
		if err != nil {
			log.Println("CCS sent error:", err)
		}
		go readHandler(m)
	}
}

func readHandler(m *ccs.InMsg) {
	log.Println("Incoming CCS message:", m)
}
