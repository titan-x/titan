package titan_test

import (
	"fmt"
	"log"

	"github.com/titan-x/titan"
)

const debug = false

// Example demonstrating the Titan server.
func Example() {
	s, err := titan.NewServer("127.0.0.1:3000")
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}
	defer s.Close()
	// if err := s.ListenAndServe(); err != nil {
	// 	log.Fatalf("error closing server: %v", err)
	// }

	fmt.Println("Server started.")

	// Output: Server started.
}
