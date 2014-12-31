package main

import (
	"log"
	"strconv"

	"github.com/nbusy/gcm/ccs"
)

var users = make(map[uint32]User)

func main() {
	c, err := ccs.Connect(Conf.GCM.CCSHost, Conf.GCM.SenderID, Conf.GCM.APIKey(), Conf.App.Debug)
	if err != nil {
		log.Fatalln("Failed to connect to GCM CCS with error:", err)
	}
	log.Println("NBusy message server started.")

	for {
		m, err := c.Receive()
		if err != nil {
			log.Println("Error receiving message:", err)
		}

		go readHandler(m)
	}
}

func readHandler(m *ccs.InMsg) {
	ids := m.Data["to_user"]
	if ids == "" {
		log.Printf("Unknown message from device: %+v\n", m)
		return
	}

	id64, err := strconv.ParseUint(ids, 10, 32)
	if err != nil || id64 == 0 {
		log.Printf("Invalid use ID specific in to_user data field in message from device: %+v\n", m)
		return
	}

	id := uint32(id64)
	user, ok := users[id]
	if !ok {
		log.Printf("User not found in user list: %+v\n", m)
	}

	user.Devices[0].Send(m.Data)
}
