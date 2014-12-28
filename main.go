package main

import (
	"log"

	"github.com/nbusy/gcm/ccs"
)

// Chat is a private or group chat.
type Chat struct {
	ID    uint64
	Users []uint64
}

var chats = make(map[string]Chat)

// user -> id (user or chat id) -> message
// delivery status -> user
// read status -> user

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
