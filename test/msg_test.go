package test

import "testing"

func TestSendEcho(t *testing.T) {
	s := NewServerHelper(t).SeedDB()
	defer s.Stop()
	c := NewClientHelper(t, s).AsUser(&s.SeedData.User1).Dial()
	defer c.Close()

	id := c.WriteRequest("msg.echo", map[string]string{"echo": "echo"})
	res := c.ReadRes(nil)

	resMap := res.Result.(map[string]interface{})
	if res.ID != id || resMap["echo"] != "echo" {
		t.Fatal("Failed to receive echo message in proper format:", res)
	}

	// t.Fatal("Failed to send a message to the echo user")
	// t.Fatal("Failed to send batch message to the echo user")
	// t.Fatal("Failed to send large message to the echo user")
	// t.Fatal("Did not receive ACK for a message sent")
	// t.Fatal("Failed to receive a response from echo user")
	// t.Fatal("Could not send an ACK for an incoming message")
}

type sendMsgReq struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type recvMsgReq struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func TestSendMsgOnline(t *testing.T) {
	s := NewServerHelper(t).SeedDB()
	defer s.Stop()
	c1 := NewClientHelper(t, s).AsUser(&s.SeedData.User1).Dial()
	defer c1.Close()
	c2 := NewClientHelper(t, s).AsUser(&s.SeedData.User2).Dial()
	defer c2.Close()

	// send msg.recv request from user 2 to announce availability and complete client-cert auth
	c2.WriteRequest("msg.recv", nil)
	res := c2.ReadRes(nil)
	if res.Result != "ACK" {
		t.Fatal("Failed to send msg.recv request from client 2 to server:", res)
	}

	// send message to user 2 with a basic hello message
	c1.WriteRequest("msg.send", sendMsgReq{To: "2", Message: "Hello, how are you?"})
	res = c1.ReadRes(nil)
	if res.Result != "ACK" {
		t.Fatal("Failed to send message to user 2:", res)
	}

	// receive the hello message from user 1 (online) as user 2 (online)
	var c2r recvMsgReq
	c2req := c2.ReadReq(&c2r)
	if c2r.From != "1" {
		t.Fatal("Received message from wrong sender instead of 1:", c2r.From)
	} else if c2r.Message != "Hello, how are you?" {
		t.Fatal("Received wrong message content:", c2r.Message)
	}

	c2.WriteResponse(c2req.ID, "ACK", nil)

	// send back a hello response to user 1 (online) as user 2 (online)
	c2.WriteRequest("msg.send", sendMsgReq{To: "1", Message: "I'm fine, thank you."})
	res = c2.ReadRes(nil)
	if res.Result != "ACK" {
		t.Fatal("Failed to send message to user 1:", res)
	}

	// receive hello response from user 1 (online) as user 2 (online)
	var c1r recvMsgReq
	c1req := c1.ReadReq(&c1r)
	if c1r.From != "2" {
		t.Fatal("Received message from wrong sender instead of 2:", c1r.From)
	} else if c1r.Message != "I'm fine, thank you." {
		t.Fatal("Received wrong message content:", c1r.Message)
	}

	c1.WriteResponse(c1req.ID, "ACK", nil)

	// todo: verify that there are no pending requests for either user 1 or 2
	// todo: below is a placeholder since writing last ACK response will never finish as we never wait for it
	c1.WriteRequest("msg.echo", map[string]string{"echo": "echo"})
	resfin := c1.ReadRes(nil).Result.(map[string]interface{})["echo"]
	if resfin != "echo" {
		t.Fatal("Last echo did return an invalid response:", resfin)
	}
}

func TestSendMsgOffline(t *testing.T) {
	s := NewServerHelper(t).SeedDB()
	defer s.Stop()
	c1 := NewClientHelper(t, s).AsUser(&s.SeedData.User1).Dial()
	defer c1.Close()

	// send message to user 2 with a basic hello message
	c1.WriteRequest("msg.send", sendMsgReq{To: "2", Message: "Hello, how are you?"})
	res := c1.ReadRes(nil)
	if res.Result != "ACK" {
		t.Fatal("Failed to send message to user 2:", res)
	}

	// connect as user 2 and send msg.recv request to announce availability and complete client-cert auth
	c2 := NewClientHelper(t, s).AsUser(&s.SeedData.User2).Dial()
	defer c2.Close()

	c2.WriteRequest("msg.recv", nil)
	res = c2.ReadRes(nil)
	if res.Result != "ACK" {
		t.Fatal("Failed to send msg.recv request from client 2 to server:", res)
	}

	// receive the hello message from user 1 (online) as user 2 (was offline at the time message was sent)
	var c2r recvMsgReq
	c2req := c2.ReadReq(&c2r)
	if c2r.From != "1" {
		t.Fatal("Received message from wrong sender instead of 1:", c2r.From)
	} else if c2r.Message != "Hello, how are you?" {
		t.Fatal("Received wrong message content:", c2r.Message)
	}

	c2.WriteResponse(c2req.ID, "ACK", nil)

	// todo: verify that there are no pending requests for either user 1 or 2
	// todo: below is a placeholder since writing last ACK response will never finish as we never wait for it
	c1.WriteRequest("msg.echo", map[string]string{"echo": "echo"})
	resfin := c1.ReadRes(nil).Result.(map[string]interface{})["echo"]
	if resfin != "echo" {
		t.Fatal("Last echo did return an invalid response:", resfin)
	}

	// todo: as client_helper is implicitly logging errors with t.Fatal(), we can't currently add useful information like below:
	// t.Fatal("Failed to receive queued messages after coming online")
	// t.Fatal("Failed to send ACK for received message queue")
}

func TestSendAsync(t *testing.T) {
	// test case to do all of the following simultaneously to test the async nature of titan server
	// - cert.auth
	// - msg.recv
	// - msg.send (bath to multiple people where some of whom are online)
}
