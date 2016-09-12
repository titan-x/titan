package test

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/neptulon/neptulon/middleware"
	"github.com/titan-x/titan/client"
	"github.com/titan-x/titan/models"
)

// ClientHelper is a Titan Client wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ClientHelper struct {
	Client *client.Client
	User   *models.User

	testing    *testing.T
	serverAddr string
	inMsgsChan chan []models.Message
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
	ch := &ClientHelper{
		Client:     c,
		testing:    t,
		serverAddr: addr,
		inMsgsChan: make(chan []models.Message),
	}
	c.MiddlewareFunc(middleware.LoggerWithPrefix("client"))
	c.InMsgHandler(ch.inMsgHandler)
	return ch
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
func (ch *ClientHelper) AsUser(u *models.User) *ClientHelper {
	ch.User = u
	return ch
}

// GoogleAuthSync is synchronous version of Client.GoogleAuth method.
// Google OAuth token is exchanged for a JWT token. If any user was assigned with AsUser, the new JWT token is stored in the user's profile.
func (ch *ClientHelper) GoogleAuthSync(oauthToken string) *ClientHelper {
	gotRes := make(chan bool)

	if err := ch.Client.GoogleAuth(oauthToken, func(jwtToken string) error {
		if jwtToken == "" {
			ch.testing.Fatalf("auth.google request failed with error: %v", "") // todo: retrieve error
		}
		ch.User.JWTToken = jwtToken
		gotRes <- true
		return nil
	}); err != nil {
		ch.testing.Fatalf("google authentication request failed: %v", err)
	}

	select {
	case <-gotRes:
	case <-time.After(time.Second * 3):
		ch.testing.Fatal("did not get an auth.jwt response in time")
	}
	return ch
}

// JWTAuthSync does JWT authentication with the token belonging the the user assigned with AsUser method.
// This method runs synchronously and blocks until authentication response is received (or connection is closed by server).
func (ch *ClientHelper) JWTAuthSync() *ClientHelper {
	gotRes := make(chan bool)

	if err := ch.Client.JWTAuth(ch.User.JWTToken, func(ack string) error {
		if ack != client.ACK {
			ch.testing.Fatalf("server did not ACK our auth.jwt request: %v", ack)
		}
		gotRes <- true
		return nil
	}); err != nil {
		ch.testing.Fatalf("jwt authentication request failed: %v", err)
	}

	select {
	case <-gotRes:
	case <-time.After(time.Second * 3):
		ch.testing.Fatal("did not get an auth.jwt response in time")
	}
	return ch
}

// EchoSync is synchronous version of Client.Echo method.
func (ch *ClientHelper) EchoSync(message string) *ClientHelper {
	gotRes := make(chan bool)

	if err := ch.Client.Echo(models.Message{Message: message}, func(msg *models.Message) error {
		if msg.Message != message {
			ch.testing.Fatalf("expected: %v, got: %v", message, msg.Message)
		}
		gotRes <- true
		return nil
	}); err != nil {
		ch.testing.Fatal(err)
	}

	select {
	case <-gotRes:
	case <-time.After(time.Second * 3):
		ch.testing.Fatal("did not get an echo response in time")
	}
	return ch
}

// SendMessagesSync is synchronous version of Client.SendMessages method.
func (ch *ClientHelper) SendMessagesSync(messages []models.Message) *ClientHelper {
	gotRes := make(chan bool)

	if err := ch.Client.SendMessages(messages, func(ack string) error {
		if ack != client.ACK {
			ch.testing.Fatalf("failed to send hello message to user %v: %v", messages[0].To, ack)
		}
		gotRes <- true
		return nil
	}); err != nil {
		ch.testing.Fatal(err)
	}

	select {
	case <-gotRes:
	case <-time.After(time.Second * 3):
		ch.testing.Fatal("did not get an msg.send response in time")
	}
	return ch
}

// GetMessagesWait waits for and returns incoming messages.
// If no message arrives within the timeout, test fails.
func (ch *ClientHelper) GetMessagesWait() []models.Message {
	select {
	case m := <-ch.inMsgsChan:
		return m
	case <-time.After(time.Second * 3):
		ch.testing.Fatal("GetMessagesWait timeout")
	}
	return nil
}

// CloseWait closes a connection.
// Waits till all the goroutines handling messages quit.
func (ch *ClientHelper) CloseWait() {
	if err := ch.Client.Close(); err != nil {
		ch.testing.Fatal("Failed to close the client connection:", err)
	}

	if os.Getenv("TRAVIS") != "" || os.Getenv("CI") != "" {
		time.Sleep(time.Millisecond * 50)
	} else {
		time.Sleep(time.Millisecond * 5)
	}
}

func (ch *ClientHelper) inMsgHandler(m []models.Message) error {
	ch.inMsgsChan <- m
	return nil
}
