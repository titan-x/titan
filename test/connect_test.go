package test

import (
	"testing"
	"time"

	"github.com/nbusy/devastator/neptulon/jsonrpc"
)

func TestClientDisconnect(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)

	closeClientConn(t, c)
	stopServer(t, s)
}

func TestServerDisconnect(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)

	stopServer(t, s)
	closeClientConn(t, c)
}

func TestClientClose(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)

	writeMsg(t, c, jsonrpc.Request{Method: "close"})

	// todo: how to wait flush gorotune spawn
	// todo: with nanosecond wait, client disconnected doesn't happen!
	time.Sleep(time.Nanosecond)

	closeClientConn(t, c)
	stopServer(t, s)
}

func TestServerClose(t *testing.T) {
	// t.Fatal("Server method.close request was not handled properly")
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
