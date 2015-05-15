package devastator_test

import (
	"log"

	"github.com/nbusy/devastator"
)

// Example demonstrating the use of Devastator server.
func Example() {
	s, err := devastator.NewServer(nil, nil, "", false)
	if err == nil && s != nil {
		log.Println("Connected")
	}
}
