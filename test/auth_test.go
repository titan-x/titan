package test

import (
	"testing"

	"github.com/nbusy/devastator/neptulon/jsonrpc"
)

func TestAuth(t *testing.T) {
	// t.Fatal("Unauthorized clients cannot call any function other than method.auth and method.close") // call to randomized and all registered routes here
	// t.Fatal("Anonymous calls to method.auth and method.close should be allowed")
}

func TestGoogleAuth(t *testing.T) {
	// t.Fatal("Google+ first sign-in (registration) failed with valid credentials")
	// t.Fatal("Google+ second sign-in (regular) failed with valid credentials")
	// t.Fatal("Google+ sign-in passed with invalid credentials")
	// t.Fatal("Authentication was not ACKed")
}

func TestClientCertAuth(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)

	writeMsg(t, c, jsonrpc.Request{ID: "123", Method: "auth.cert"}) // should be a variadic fn(method, params...)
	m := readMsg(t, c)

	if m.ID != "123" && m.Result != "ACK" {
		t.Fatal("Authentication failed with a valid client certificate. Got server response:", m)
	}

	closeClientConn(t, c)
	stopServer(t, s)

	// t.Fatal("Authenticated with invalid/expired client certificate")
	// t.Fatal("Authentication was not ACKed")
}

func TestTokenAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid token")
	// t.Fatal("Authenticated with invalid/expired token")
	// t.Fatal("Authentication was not ACKed")
}
