package test

import "testing"

func TestClientDisconnect(t *testing.T) {
	// todo: we need to verify that events occur in the order that we want them (either via event hooks or log analysis)
	// this seems like a listener test than a integration test from a client perspective
	s := getServer(t)
	c := getClientConnWithClientCert(t)
	if err := c.Close(); err != nil {
		t.Fatal("Failed to close the client connection:", err)
	}
	if err := s.Stop(); err != nil {
		t.Fatal("Failed to stop the server:", err)
	}
	wg.Wait()
}

func TestClientClose(t *testing.T) {
	// t.Fatal("Client method.close request was not handled properly")
}

func TestSendClose(t *testing.T) {
	// t.Fatal("Server method.close request was not handled properly")
}

func TestServerDisconnect(t *testing.T) {
	// t.Fatal("Server disconnect was not handled gracefully")
}

func TestServerClose(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)
	if err := s.Stop(); err != nil {
		t.Fatal("Failed to stop the server:", err)
	}
	if err := c.Close(); err != nil {
		t.Fatal("Failed to close the client connection:", err)
	}
	wg.Wait()

	// test what happens when there are outstanding connections and/or requests that are being handled
	// destroying queues and other stuff during Close() might cause existing request handles to malfunction
}

func TestStop(t *testing.T) {
	// todo: this should be a listener/queue test if we don't use any goroutines in the Server struct methods
	// t.Fatal("Failed to stop the server: not all the goroutines were terminated properly")
	// t.Fatal("Failed to stop the server: server did not wait for ongoing read/write operations")
}

func TestConnTimeout(t *testing.T) {
	// t.Fatal("Send timout did not occur")
	// t.Fatal("Wait timeout did not occur")
	// t.Fatal("Read timeout did not occur")
}
