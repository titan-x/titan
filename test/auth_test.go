package test

import (
	"os"
	"testing"

	"github.com/nbusy/devastator"
)

func TestAuth(t *testing.T) {
	// t.Fatal("Unauthorized clients cannot call any function other than method.auth and method.close") // call to randomized and all registered routes here
	// t.Fatal("Anonymous calls to method.auth and method.close should be allowed")
}

func TestValidClientCertAuth(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).DefaultCert().Dial()
	defer c.Close()
	id := c.WriteRequest("echo", nil)
	_, res, _ := c.ReadMsg(nil)

	if res.ID != id {
		t.Fatal("Authentication failed with a valid client certificate. Got server response:", res)
	}
}

// todo: no cert, no signature cert, invalid CA signed cert, expired cert...
func TestInvalidClientCertAuth(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).Dial()
	defer c.Close()

	_ = c.WriteRequest("echo", nil)

	if !c.VerifyConnClosed() {
		t.Fatal("Authenticated successfully with invalid client certificate.")
	}
}

func TestGoogleAuth(t *testing.T) {
	token := os.Getenv("GOOGLE_ACCESS_TOKEN")
	if token == "" {
		t.Skip("Missing 'GOOGLE_ACCESS_TOKEN' environment variable. Skipping Google sign-in testing.")
	}

	s := NewServerHelper(t)
	c := NewClientHelper(t).Dial()

	c.WriteRequest("auth.google", map[string]string{"accessToken": token})
	res := c.ReadRes(&devastator.CertResponse{}) // todo: we need to be able to specify return type here, otherwise we get a map[]

	if res.Error != nil {
		t.Fatal("Google+ first sign-in/registration failed with valid credentials:", res.Error)
	}

	c.Close()
	s.Stop()

	// now connect to server with our new client certificate
	r := res.Result.(*devastator.CertResponse)
	cert, key := r.Cert, r.Key

	s = NewServerHelper(t)
	c = NewClientHelper(t).Cert(cert, key).Dial()

	_ = c.WriteRequest("echo", nil)
	res = c.ReadRes(nil)

	if res.Error != nil {
		t.Fatal("Failed to connect to the server with certificates created after Google+ sign-in:", res.Error)
	}

	c.Close()
	s.Stop()
}

func TestInvalidGoogleAuth(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).Dial()
	defer c.Close()

	// t.Fatal("Google+ second sign-in (regular) failed with valid credentials")
	// t.Fatal("Google+ sign-in passed with invalid credentials")
}
