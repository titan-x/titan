package test

import (
	"sync"
	"testing"

	"github.com/titan-x/titan/client"
)

func TestAuth(t *testing.T) {
	// t.Fatal("Unauthorized clients cannot call any function other than method.auth and method.close") // call to randomized and all registered routes here
	// t.Fatal("Anonymous calls to method.auth and method.close should be allowed")
}

func TestValidToken(t *testing.T) {
	sh := NewServerHelper(t).SeedDB()
	defer sh.ListenAndServe().CloseWait()

	ch := sh.GetClientHelper().AsUser(&sh.SeedData.User1)
	defer ch.Connect().CloseWait()

	var wg sync.WaitGroup
	wg.Add(1)
	msg := "Lorem ip sum"

	ch.Client.Echo(map[string]string{"message": msg, "token": sh.SeedData.User1.JWTToken}, func(m *client.Message) error {
		defer wg.Done()
		if m.Message != msg {
			t.Fatalf("expected: %v, got: %v", "wow", m.Message)
		}
		return nil
	})

	wg.Wait()
}

func TestInvalidToken(t *testing.T) {
	sh := NewServerHelper(t).SeedDB()
	defer sh.ListenAndServe().CloseWait()

	ch := sh.GetClientHelper().AsUser(&sh.SeedData.User1)
	defer ch.Connect().CloseWait()

	var wg sync.WaitGroup
	wg.Add(1)
	msg := "Lorem ip sum"

	ch.Client.DisconnHandler(func(c *client.Client) {
		wg.Done()
	})

	ch.Client.Echo(map[string]string{"message": msg, "token": "abc-invalid-token-!"}, func(m *client.Message) error {
		t.Fatal("authenticated with invalid token")
		return nil
	})

	wg.Wait()

	// todo: no token, un-signed token, invalid token signature, expired token...
}

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
