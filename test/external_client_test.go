package test

import (
	"flag"
	"testing"
)

var ext = flag.Bool("ext", false, "Run external client test case.")

// Helper method for testing client implementations in other languages.
// Flow of events for this function is:
// * Send a {"method":"echo", "params":{"message": "..."}} request to client upon first 'echo' request from client,
//   and verify that message body is echoed properly in the response body.
// * Echo any incoming request message body as is within a response message.
// * Repeat ad infinitum, until {"method":"close", "params":"{"message": "..."}"} is received. Close message body is logged.
func TestExternalClient(t *testing.T) {
	sh := NewServerHelper(t).SeedDB()
	defer sh.ListenAndServe().CloseWait()
	// m := "Hello from Neptulon server!"

}
