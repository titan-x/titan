package titan_test

import (
	"log"

	"github.com/titan-x/titan"
)

const debug = false

// Example demonstrating the Titan server.
func Example() {
	s, err := titan.NewServer("127.0.0.1:3000")
	if err != nil {
		log.Fatalln("Errored while creating a new server instance:", err)
	}

	if s != nil {
		log.Println("Connected")

	}

	// ** Output: Server started
}
