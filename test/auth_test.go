package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/titan-x/titan/client"
	"github.com/titan-x/titan/models"
)

func TestAuth(t *testing.T) {
	// t.Fatal("Unauthorized clients cannot call any function other than method.auth and method.close") // call to randomized and all registered routes here
	// t.Fatal("Anonymous calls to method.auth and method.close should be allowed")
}

func TestValidToken(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	ch := sh.GetClientHelper().AsUser(&sh.SeedData.User1).Connect().JWTAuthSync()
	defer ch.CloseWait()

	ch.EchoSync("Ola!")
}

func TestInvalidToken(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	ch := sh.GetClientHelper().Connect()
	defer ch.CloseWait()

	gotMsg, closed := make(chan bool), make(chan bool)
	ch.Client.DisconnHandler(func(c *client.Client) {
		closed <- true
	})
	ch.Client.Echo(map[string]string{"message": "Lorem ip sum", "token": "abc-invalid-token-!"}, func(m *models.Message) error {
		gotMsg <- true
		return nil
	})

	select {
	case <-gotMsg:
		t.Fatal("authenticated with invalid token")
	case <-closed:
		log.Println("test: server closed connection as expected")
	case <-time.After(time.Second):
	}

	// todo: no token, un-signed token, invalid token signature, expired token...
}

type googleAuthRes struct {
	Cert, Key []byte
}

func TestGoogleAuth(t *testing.T) {
	token := os.Getenv("GOOGLE_ACCESS_TOKEN")
	if token == "" {
		t.Skip("Missing 'GOOGLE_ACCESS_TOKEN' environment variable. Skipping Google sign-in testing.")
	}

	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	// authenticate with Google OAuth token and get JWT token
	ch := sh.GetClientHelper().AsUser(&sh.SeedData.User1).Connect().GoogleAuthSync(token)

	// send an echo message to validate that we are authenticated properly
	ch.EchoSync("testing echo message after google auth")
	ch.CloseWait()

	// now connect to server with our new JWT token auto assigned by Google auth helper function
	ch.Connect().JWTAuthSync()
	ch.SendMessagesSync([]models.Message{models.Message{To: "2", Message: "Hi!"}})

	ch.CloseWait()
}

func TestInvalidGoogleAuth(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	// authenticate with Google OAuth token and get JWT token
	// ch := sh.GetClientHelper().AsUser(&sh.SeedData.User1).Connect().GoogleAuthSync(token)

	// t.Fatal("Google+ second sign-in (regular) failed with valid credentials")
	// t.Fatal("Google+ sign-in passed with invalid credentials")
}
