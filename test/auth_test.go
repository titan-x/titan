package test

import "testing"

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

func TestValidClientCertAuth(t *testing.T) {
	h := NewClientServerHelper(t, true)
	defer h.Close()

	id := h.Client.WriteRequest("auth.cert", nil)
	m := h.Client.ReadMsg()

	if m.ID != id || m.Result != "OK" {
		t.Fatal("Authentication failed with a valid client certificate. Got server response:", m)
	}
}

func TestInvalidClientCertAuth(t *testing.T) {
	h := NewClientServerHelper(t, false)
	defer h.Close()

	id := h.Client.WriteRequest("auth.cert", nil)
	m := h.Client.ReadMsg()

	if m.ID != id || m.Result == "OK" {
		t.Fatal("Authenticated with invalid/expired client certificate. Got server response:", m)
	}
}

func TestTokenAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid token")
	// t.Fatal("Authenticated with invalid/expired token")
	// t.Fatal("Authentication was not ACKed")
}
