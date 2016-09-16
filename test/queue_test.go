package test

import (
	"testing"

	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/models"
)

func TestToCaseInsensitive(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	ch := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect().JWTAuthSync()
	defer ch.CloseWait()

	ch.SendMessagesSync([]models.Message{models.Message{To: "Echo", Message: "Ola!"}}) // Echo capital on purpose to test

	msgs := ch.GetMessagesWait()

	msg := msgs[0]
	if msg.From != "echo" {
		t.Fatalf("expected message from: echo, got: %v", msg.From)
	}
	if msg.Message != "Ola!" {
		t.Fatalf("expected message from: Ola!, got: %v", msg.Message)
	}
}

func TestMultipleConcurrentQueuedMessages(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	ch := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect().JWTAuthSync()
	defer ch.CloseWait()

	ch.SendMessagesSync([]models.Message{
		models.Message{To: "echo", Message: "message-1"},
		models.Message{To: "echo", Message: "message-2"},
		models.Message{To: "echo", Message: "message-3"},
		models.Message{To: "echo", Message: "message-4"}})

	// todo: do this in a for loop in a go-routine to be concurrent-realistic
	ch.GetMessagesWait()

	// if len(msgs) != 4 {
	// 	t.Fatalf("expected 4 messages, got %v", len(msgs))
	// }
}
