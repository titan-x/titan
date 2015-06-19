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

func TestGoogleRegister(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).Dial()
	defer c.Close()

	// h.Client.WriteRequest("auth.google", map[string]string{"OAuthToken": "1234"})
	// m := h.Client.ReadMsg()

	// should get client cert and try connecting with it again
}

func TestValidClientCertAuth(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).DefaultCert().Dial()
	defer c.Close()

	id := c.WriteRequest("auth.cert", nil)
	m := c.ReadMsg()

	if m.ID != id || m.Result != "OK" {
		t.Fatal("Authentication failed with a valid client certificate. Got server response:", m)
	}
}

func TestInvalidClientCertAuth(t *testing.T) {
	s := NewServerHelper(t)
	defer s.Stop()
	c := NewClientHelper(t).Dial()
	defer c.Close()

	_ = c.WriteRequest("auth.cert", nil)
	m := c.ReadMsg()

	if m.Result != nil || m.Error.Code != 666 || m.Error.Message != "Invalid client certificate." {
		t.Fatal("Authenticated successfully with invalid client certificate. Got server response:", m)
	}
}

func TestTokenAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid token")
	// t.Fatal("Authenticated with invalid/expired token")
	// t.Fatal("Authentication was not ACKed")
}
