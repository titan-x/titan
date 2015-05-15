package devastator_test

import (
	"log"

	"github.com/nbusy/devastator"
)

// todo: add server.go setup example

// Example demonstrating the use of Devastator server.
func Example() {
	s, err := devastator.NewServer(nil, nil, "", false)
	if err == nil && s != nil {
		log.Println("Connected")
	}
}
