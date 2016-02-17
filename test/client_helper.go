package test

import (
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/titan-x/titan"
	"github.com/titan-x/titan/client"
)

// ClientHelper is a Titan Client wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ClientHelper struct {
	Client *client.Client
	User   *titan.User

	testing    *testing.T
	serverAddr string
}

// NewClientHelper creates a new client helper object.
func NewClientHelper(t *testing.T, addr string) *ClientHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode.")
	}

	c, err := client.NewClient()
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
	ch.User = u
	return ch
}

// JWTAuth does JWT authentication with the token belonging the the user assigned with AsUser method.
// This method runs synchronously and blocks until authentication response is received (or connection is closed by server).
func (ch *ClientHelper) JWTAuth() *ClientHelper {
	var wg sync.WaitGroup
	wg.Add(1)
	timer := time.AfterFunc(time.Millisecond*100, func() {
		wg.Done()
	})

	if err := ch.Client.JWTAuth(ch.User.JWTToken, func(ack string) error {
		timer.Stop()
		defer wg.Done()
		if ack != "ACK" {
			return fmt.Errorf("server did not ACK our auth.jwt request: %v", ack)
		}
		return nil
	}); err != nil {
		timer.Stop()
		defer wg.Done()
		ch.testing.Fatalf("authentication failed: %v", err)
	}

	wg.Wait()
	return ch
}

// CloseWait closes a connection.
// Waits till all the goroutines handling messages quit.
func (ch *ClientHelper) CloseWait() {
	if err := ch.Client.Close(); err != nil {
		ch.testing.Fatal("Failed to close the client connection:", err)
	}
	time.Sleep(time.Millisecond * 5)
	if os.Getenv("TRAVIS") != "" || os.Getenv("CI") == "" {
		time.Sleep(time.Millisecond * 50)
	}
}
