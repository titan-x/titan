package neptulon_test

import (
	"log"

	"github.com/nbusy/devastator"
)

const debug = false

// Example demonstrating the Devastator server.
func Example() {
	s, err := devastator.NewServer(nil, nil, "", debug)
	if err != nil {
		log.Fatalln("Errored while creating a new server instance:", err)
	}

	if s != nil {
		log.Println("Connected")

	}

	// ** Output: Server started
}
