package test

import (
	"testing"

	"github.com/nbusy/devastator"
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

	// test what happens when there are outstanding connections and/or requests that are being handled
	// destroying queues and other stuff during Close() might cause existing request handles to malfunction

	stopServer(t, s)
	closeClientConn(t, c)
}

func TestClientClose(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)

	writeMsg(t, c, devastator.ReqMsg{Method: "close"})

	closeClientConn(t, c)
	stopServer(t, s)
}

func TestServerClose(t *testing.T) {
	// t.Fatal("Server method.close request was not handled properly")
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
