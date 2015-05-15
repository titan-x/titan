package devastator_test

import (
	"log"

	"github.com/nbusy/devastator"
)

const debug = false

// Example demonstrating the use of Devastator server.
func Example() {
	s, err := devastator.NewServer(nil, nil, "", debug)
	if err == nil && s != nil {
		log.Println("Connected")
	}
}
