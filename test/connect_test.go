package test

import (
	"testing"
	"time"

	"github.com/titan-x/titan/data"
)

func TestClientDisconnect(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	ch := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect()

	time.Sleep(time.Millisecond * 10)

	ch.CloseWait()
	sh.CloseWait()

	// todo: validate log output order
}

func TestServerDisconnect(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	ch := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect()

	time.Sleep(time.Millisecond * 10)

	sh.CloseWait()
	ch.CloseWait()

	// todo: validate log output order
}

func TestClientClose(t *testing.T) {
	sh := NewServerHelper(t).ListenAndServe()
	defer sh.CloseWait()

	ch := sh.GetClientHelper().AsUser(&data.SeedUser1).Connect()
	defer ch.CloseWait()

	// c.WriteNotification("conn.close", nil)
	// todo: verify that connection is closed by listner before client does

	// note: with nanosecond wait, client disconnected doesn't happen because we close client and the server faster than
	// listener.handleMsg() goroutine can be spawned as we're never waiting for an answer
	// graceful listener closing attempts also fail as reqWG is never incremented in the same method
	// time.Sleep(time.Nanosecond)
}

func TestServerClose(t *testing.T) {
	// t.Fatal("Server->client method.close request was not handled properly")
	// t.Fatal("ACK for Server->client method.close request was not received")
}

func TestMultiConn(t *testing.T) {
	// t.Fatal("Failed to handle randomly connecting disconnecting users")
}

func TestStop(t *testing.T) {
	// todo: this should be a listener/queue test if we don't use any goroutines in the Server struct methods
	// t.Fatal("Failed to stop the server: not all the goroutines were terminated properly")
	// t.Fatal("Failed to stop the server: server did not wait for ongoing read/write operations")
	// t.Fatal("Server did not release port 3001 after closing.")
}

func TestConnTimeout(t *testing.T) {
	// t.Fatal("Send timout did not occur")
	// t.Fatal("Wait timeout did not occur")
	// t.Fatal("Read timeout did not occur")
}
