package test

import "testing"

func TestAuth(t *testing.T) {
	// t.Fatal("Unauthorized clients cannot call any function other than method.auth and method.close") // call to randomized and all registered routes here
	// t.Fatal("Anonymous calls to method.auth and method.close should be allowed")
}

func TestValidClientCertAuth(t *testing.T) {

	// todo: add s.Start() to be consistent with inner server helpers

	s := NewServerHelper(t).SeedDB()
	defer s.Stop()

	c := s.GetClientHelper().AsUser(&s.SeedData.User1).Connect()
	defer c.Close()

	// id := c.WriteRequest("msg.echo", nil)
	// res := c.ReadRes(nil)
	//
	// if res.ID != id {
	// 	t.Fatal("Authentication failed with a valid client certificate. Got server response:", res)
	// }
}

//
// func TestInvalidClientCertAuth(t *testing.T) {
// 	s := NewServerHelper(t)
// 	defer s.Stop()
// 	c := NewConnHelper(t, s).Dial()
// 	defer c.Close()
//
// 	_ = c.WriteRequest("msg.echo", nil)
//
// 	if !c.VerifyConnClosed() {
// 		t.Fatal("Authenticated successfully with invalid client certificate.")
// 	}
//
// 	// todo: no cert, no signature cert, invalid CA signed cert, expired cert...
// }
//
// type googleAuthRes struct {
// 	Cert, Key []byte
// }
//
// func TestGoogleAuth(t *testing.T) {
// 	token := os.Getenv("GOOGLE_ACCESS_TOKEN")
// 	if token == "" {
// 		t.Skip("Missing 'GOOGLE_ACCESS_TOKEN' environment variable. Skipping Google sign-in testing.")
// 	}
//
// 	s := NewServerHelper(t)
// 	c := NewConnHelper(t, s).Dial()
//
// 	c.WriteRequest("auth.google", map[string]string{"accessToken": token})
// 	var resData googleAuthRes
// 	res := c.ReadRes(&resData)
//
// 	if res.Error != nil {
// 		t.Fatal("Google+ first sign-in/registration failed with valid credentials:", res.Error)
// 	}
//
// 	c.Close()
// 	s.Stop()
//
// 	// now connect to server with our new client certificate
// 	s = NewServerHelper(t)
// 	c = NewConnHelper(t, s).WithCert(resData.Cert, resData.Key).Dial()
//
// 	_ = c.WriteRequest("msg.echo", nil)
// 	res = c.ReadRes(nil)
//
// 	if res.Error != nil {
// 		t.Fatal("Failed to connect to the server with certificates created after Google+ sign-in:", res.Error)
// 	}
//
// 	c.Close()
// 	s.Stop()
// }
//
// func TestInvalidGoogleAuth(t *testing.T) {
// 	s := NewServerHelper(t)
// 	defer s.Stop()
// 	c := NewConnHelper(t, s).Dial()
// 	defer c.Close()
//
// 	// t.Fatal("Google+ second sign-in (regular) failed with valid credentials")
// 	// t.Fatal("Google+ sign-in passed with invalid credentials")
// }
