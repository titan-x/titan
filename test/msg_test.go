package test

import (
	"sync"
	"testing"

	"github.com/titan-x/titan/client"
)

func TestSendEcho(t *testing.T) {
	sh := NewServerHelper(t).SeedDB()
	defer sh.ListenAndServe().CloseWait()

	ch := sh.GetClientHelper().AsUser(&sh.SeedData.User1)
	defer ch.Connect().CloseWait()

	var wg sync.WaitGroup
	wg.Add(1)
	m := "Ola!"
	ch.Client.Echo(map[string]string{"message": m, "token": sh.SeedData.User1.JWTToken}, func(msg *client.Message) error {
		defer wg.Done()
		if msg.Message != m {
			t.Fatalf("expected: %v, got: %v", m, msg.Message)
		}
		return nil
	})

	wg.Wait()

	// t.Fatal("Failed to send a message to the echo user")
	// t.Fatal("Failed to send batch message to the echo user")
	// t.Fatal("Failed to send large message to the echo user")
	// t.Fatal("Did not receive ACK for a message sent")
	// t.Fatal("Failed to receive a response from echo user")
	// t.Fatal("Could not send an ACK for an incoming message")
}

func TestSendMsgOnline(t *testing.T) {
	sh := NewServerHelper(t).SeedDB()
	defer sh.ListenAndServe().CloseWait()

	ch1 := sh.GetClientHelper().AsUser(&sh.SeedData.User1)
	defer ch1.Connect().JWTAuth().CloseWait()

	ch2 := sh.GetClientHelper().AsUser(&sh.SeedData.User2)
	defer ch2.Connect().JWTAuth().CloseWait()

	var wg sync.WaitGroup

	// send a hello message from user 1 to user 2
	wg.Add(1)
	if err := ch1.Client.SendMessages([]client.Message{client.Message{To: "2", Message: "Hello, how are you?"}}, func(ack string) error {
		defer wg.Done()
		if ack != "ACK" {
			t.Fatal("failed to send hello message to user 2:", ack)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// todo: all of SendMessages etc. now return error so handle them as above

	//
	// // receive the hello message from user 1 (online) as user 2 (online)
	// var c2r recvMsgReq
	// c2req := c2.ReadReq(&c2r)
	// if c2r.From != "1" {
	// 	t.Fatal("Received message from wrong sender instead of 1:", c2r.From)
	// } else if c2r.Message != "Hello, how are you?" {
	// 	t.Fatal("Received wrong message content:", c2r.Message)
	// }
	//
	// c2.WriteResponse(c2req.ID, "ACK", nil)
	//
	// // send back a hello response to user 1 (online) as user 2 (online)
	// c2.WriteRequest("msg.send", sendMsgReq{To: "1", Message: "I'm fine, thank you."})
	// res = c2.ReadRes(nil)
	// if res.Result != "ACK" {
	// 	t.Fatal("Failed to send message to user 1:", res)
	// }
	//
	// // receive hello response from user 1 (online) as user 2 (online)
	// var c1r recvMsgReq
	// c1req := c1.ReadReq(&c1r)
	// if c1r.From != "2" {
	// 	t.Fatal("Received message from wrong sender instead of 2:", c1r.From)
	// } else if c1r.Message != "I'm fine, thank you." {
	// 	t.Fatal("Received wrong message content:", c1r.Message)
	// }
	//
	// c1.WriteResponse(c1req.ID, "ACK", nil)
	//
	// // todo: verify that there are no pending requests for either user 1 or 2
	// // todo: below is a placeholder since writing last ACK response will never finish as we never wait for it
	// c1.WriteRequest("msg.echo", map[string]string{"echo": "echo"})
	// resfin := c1.ReadRes(nil).Result.(map[string]interface{})["echo"]
	// if resfin != "echo" {
	// 	t.Fatal("Last echo did return an invalid response:", resfin)
	// }

	wg.Wait()
}

//
// func TestSendMsgOffline(t *testing.T) {
// 	s := NewServerHelper(t).SeedDB()
// 	defer s.Stop()
// 	c1 := NewConnHelper(t, s).AsUser(&s.SeedData.User1).Dial()
// 	defer c1.Close()
//
// 	// send message to user 2 with a basic hello message
// 	c1.WriteRequest("msg.send", sendMsgReq{To: "2", Message: "Hello, how are you?"})
// 	res := c1.ReadRes(nil)
// 	if res.Result != "ACK" {
// 		t.Fatal("Failed to send message to user 2:", res)
// 	}
//
// 	// connect as user 2 and send msg.recv request to announce availability and complete client-cert auth
// 	c2 := NewConnHelper(t, s).AsUser(&s.SeedData.User2).Dial()
// 	defer c2.Close()
//
// 	c2.WriteRequest("msg.recv", nil)
// 	res = c2.ReadRes(nil)
// 	if res.Result != "ACK" {
// 		t.Fatal("Failed to send msg.recv request from client 2 to server:", res)
// 	}
//
// 	// receive the hello message from user 1 (online) as user 2 (was offline at the time message was sent)
// 	var c2r recvMsgReq
// 	c2req := c2.ReadReq(&c2r)
// 	if c2r.From != "1" {
// 		t.Fatal("Received message from wrong sender instead of 1:", c2r.From)
// 	} else if c2r.Message != "Hello, how are you?" {
// 		t.Fatal("Received wrong message content:", c2r.Message)
// 	}
//
// 	c2.WriteResponse(c2req.ID, "ACK", nil)
//
// 	// todo: verify that there are no pending requests for either user 1 or 2
// 	// todo: below is a placeholder since writing last ACK response will never finish as we never wait for it
// 	c1.WriteRequest("msg.echo", map[string]string{"echo": "echo"})
// 	resfin := c1.ReadRes(nil).Result.(map[string]interface{})["echo"]
// 	if resfin != "echo" {
// 		t.Fatal("Last echo did return an invalid response:", resfin)
// 	}
//
// 	// todo: as client_helper is implicitly logging errors with t.Fatal(), we can't currently add useful information like below:
// 	// t.Fatal("Failed to receive queued messages after coming online")
// 	// t.Fatal("Failed to send ACK for received message queue")
// }
//
// func TestSendAsync(t *testing.T) {
// 	// test case to do all of the following simultaneously to test the async nature of titan server
// 	// - cert.auth
// 	// - msg.recv
// 	// - msg.send (bath to multiple people where some of whom are online)
// }
