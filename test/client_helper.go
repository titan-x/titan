package test

import (
	"net"
	"testing"
	"time"

	"github.com/neptulon/neptulon/test"
	"github.com/titan-x/titan"
)

// ClientHelper is a Titan Client wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ClientHelper struct {
	Client *titan.Client

	ch         *test.ConnHelper // inner Neptulon Conn helper
	testing    *testing.T
	serverAddr string
	user       *titan.User
}

// NewClientHelper creates a new client helper object.
func NewClientHelper(t *testing.T, addr string) *ClientHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode.")
	}

	c, err := titan.NewClient()
	if err != nil {
		t.Fatal("Failed to create client:", err)
	}

	c.SetDeadline(10)

	return &ClientHelper{
		Client:     c,
		testing:    t,
		serverAddr: addr,
	}
}

// Connect connects to a server.
func (ch *ClientHelper) Connect() *ClientHelper {
	// retry connect in case we're operating on a very slow machine
	for i := 0; i <= 5; i++ {
		if err := ch.Client.Connect(ch.serverAddr); err != nil {
			if operr, ok := err.(*net.OpError); ok && operr.Op == "dial" && operr.Err.Error() == "connection refused" {
				time.Sleep(time.Millisecond * 50)
				continue
			} else if i == 5 {
				ch.testing.Fatalf("Cannot connect to server address %v after 5 retries, with error: %v", ch.serverAddr, err)
			}
			ch.testing.Fatalf("Cannot connect to server address %v with error: %v", ch.serverAddr, err)
		}

		if i != 0 {
			ch.testing.Logf("WARNING: it took %v retries to connect to the server, which might indicate code issues or slow machine.", i)
		}

		break
	}

	return ch
}

// AsUser attaches given user's client certificate and private key to the connection.
func (ch *ClientHelper) AsUser(u *titan.User) *ClientHelper {
	ch.user = u
	return ch
}
