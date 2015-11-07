package titan_test

import (
	"log"

	"github.com/nb-titan/titan"
)

const debug = false

// Example demonstrating the Titan server.
func Example() {
	s, err := titan.NewServer(nil, nil, nil, nil, "", debug)
	if err != nil {
		log.Fatalln("Errored while creating a new server instance:", err)
	}

	if s != nil {
		log.Println("Connected")

	}

	// ** Output: Server started
}
