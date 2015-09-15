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

func TestMsgSend(t *testing.T) {
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
		t.Fatal("Failed to send msg.recv request from client 2 with error:", res)
	}

	type sendMsgReq struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}

	c1.WriteRequest("msg.send", sendMsgReq{To: "2", Message: "lorem ip sum dolor sit amet"})
	res = c1.ReadRes(nil)
	if res.Result != "ACK" {
		t.Fatal("Failed to send message to peer with error:", res)
	}

	type recvMsgReq struct {
		From    string `json:"from"`
		Message string `json:"message"`
	}

	var r recvMsgReq
	c2.ReadReq(&r)
	if r.From != "1" {
		t.Fatal("Received message from wrong sender instead of 1:", r.From)
	} else if r.Message != "lorem ip sum dolor sit amet" {
		t.Fatal("Received wrong message content:", r.Message)
	}

	// todo: send back an answer to client 1
	// todo: verify that we receive ACK for sent response
	// todo: verify that we receive message on client 1
}

func TestMsgRecv(t *testing.T) {
	s := NewServerHelper(t).SeedDB()
	defer s.Stop()
	c1 := NewClientHelper(t, s).AsUser(&s.SeedData.User1).Dial()
	defer c1.Close()

	_ = c1.WriteRequest("msg.send", map[string]string{"to": "2", "msg": "How do you do?"})
	res := c1.ReadRes(nil)
	if res.Result != "ACK" {
		t.Fatal("Failed to send message to peer with response:", res)
	}

	c2 := NewClientHelper(t, s).AsUser(&s.SeedData.User2).Dial()
	defer c2.Close()

	// _ = c1.WriteRequest("msg.recv", nil)

	// t.Fatal("Failed to receive queued messages after coming online")
	// t.Fatal("Failed to send ACK for received message queue")
}
