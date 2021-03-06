package test

import (
	"testing"
	"time"

	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/models"
)

func TestSendEcho(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	ch := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect()
	defer ch.CloseWait()

	// not using ClientHelper.EchoSafeSync to differentiate this test from auth_test.TestValidToken
	gotRes := make(chan bool)
	m := "Ola!"
	if err := ch.Client.Echo(map[string]string{"message": m, "token": data.SeedUser1.JWTToken}, func(msg *models.Message) error {
		if msg.Message != m {
			t.Fatalf("expected: %v, got: %v", m, msg.Message)
		}
		gotRes <- true
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	select {
	case <-gotRes:
	case <-time.After(time.Second * 3):
		t.Fatal("did not get echo response in time")
	}

	// t.Fatal("Failed to send a message to the echo user")
	// t.Fatal("Failed to send batch message to the echo user")
	// t.Fatal("Failed to send large message to the echo user")
	// t.Fatal("Did not receive ACK for a message sent")
	// t.Fatal("Failed to receive a response from echo user")
	// t.Fatal("Could not send an ACK for an incoming message")
}

func TestSendMsgOnline(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	// get both user 1 and user 2 online
	ch1 := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect().JWTAuthSync()
	defer ch1.CloseWait()
	ch2 := sh.GetClientHelper().AsUser(&data.SeedUser2).Connect().JWTAuthSync()
	defer ch2.CloseWait()

	// send a hello message from user 1
	m := "Hello, how are you?"
	ch1.SendMessagesSync([]models.Message{models.Message{To: "2", Message: m}})

	// receive the hello message as user 2
	msgs := ch2.GetMessagesWait()
	if len(msgs) != 1 {
		t.Fatalf("expected message count: 1, got: %v", len(msgs))
	}
	msg := msgs[0]
	if msg.From != "1" {
		t.Fatalf("expected message from: 1, got: %v", msg.From)
	}
	if msg.Message != m {
		t.Fatalf("expected message body: %v, got: %v", m, msg.Message)
	}

	// send back a hello response from user 2
	m = "I'm fine, thank you."
	ch2.SendMessagesSync([]models.Message{models.Message{To: "1", Message: m}})

	// receive the hello response as user 1
	msgs = ch1.GetMessagesWait()
	if len(msgs) != 1 {
		t.Fatalf("expected message count: 1, got: %v", len(msgs))
	}
	msg = msgs[0]
	if msg.From != "2" {
		t.Fatalf("expected message from: 2, got: %v", msg)
	}
	if msg.Message != m {
		t.Fatalf("expected message body: %v, got: %v", m, msg.Message)
	}

	// todo: verify that there are no pending requests for either user 1 or 2
}

func TestSendMsgOffline(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	// get only user 1 online
	ch1 := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect().JWTAuthSync()
	defer ch1.CloseWait()
	ch2 := sh.GetClientHelper().AsUser(&data.SeedUser2)

	// send a hello message from user 1
	m := "Hello, how are you?"
	ch1.SendMessagesSync([]models.Message{models.Message{To: "2", Message: m}})

	// get user 2 online receive the pending hello message
	ch2.Connect().JWTAuthSync()
	defer ch2.CloseWait()
	msgs := ch2.GetMessagesWait()
	if len(msgs) != 1 {
		t.Fatalf("expected message count: 1, got: %v", len(msgs))
	}
	msg := msgs[0]
	if msg.From != "1" {
		t.Fatalf("expected message from: 1, got: %v", msg)
	}
	if msg.Message != m {
		t.Fatalf("expected message body: %v, got: %v", m, msg.Message)
	}

	// todo: verify that there are no pending requests for either user 1 or 2
}

func TestSendAsync(t *testing.T) {
	// test case to do all of the following simultaneously to test the async nature of titan server
	// - cert.auth
	// - msg.recv
	// - msg.send (bath to multiple people where some of whom are online)
}
