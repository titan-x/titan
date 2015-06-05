package neptulon

import (
	"crypto/tls"
	"crypto/x509"
	"strconv"
	"sync"
	"testing"
)

func TestLen(t *testing.T) {
	a, _ := strconv.Atoi("12344324")
	t.Log(a)
}

// todo: if we are going to expose raw Listener, this should be in integration tests, otherwise Listener should be private
func TestListener(t *testing.T) {
	msg1 := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	msg2 := "In sit amet lectus felis, at pellentesque turpis."
	msg3 := "Nunc urna enim, cursus varius aliquet ac, imperdiet eget tellus."
	msg4 := randString(45000)   //0.45 MB
	msg5 := randString(5000000) //5.0 MB

	host := "localhost:3010"
	cert, privKey, _ := genCert("localhost", 0, nil, nil, 512, "localhost", "devastator")
	l, err := Listen(cert, privKey, host, false)
	if err != nil {
		t.Fatal(err)
	}

	var listenerWG sync.WaitGroup
	listenerWG.Add(1)
	go func() {
		defer listenerWG.Done()
		l.Accept(func(conn *Conn) {},
			func(conn *Conn, msg []byte) {
				certs := conn.ConnectionState().PeerCertificates
				if len(certs) > 0 {
					t.Logf("Client connected with client certificate subject: %v\n", certs[0].Subject)
				}

				m := string(msg)
				if m != msg1 && m != msg2 && m != msg3 && m != msg4 && m != msg5 {
					t.Fatal("Sent and incoming message did not match for message:", m)
				}
			}, func(conn *Conn) {})
	}()

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		panic("failed to parse root certificate")
	}

	tlsConf := &tls.Config{RootCAs: roots}
	conn, err := tls.Dial("tcp", host, tlsConf)
	if err != nil {
		t.Fatal(err)
	}

	newconn, _ := NewConn(conn, 0, 0, 0, false)

	send(t, newconn, "ping")
	send(t, newconn, msg1)
	send(t, newconn, msg1)
	send(t, newconn, msg2)
	send(t, newconn, msg3)
	send(t, newconn, msg4)
	send(t, newconn, msg1)
	send(t, newconn, msg5)
	send(t, newconn, msg1)
	send(t, newconn, "close")

	l.reqWG.Wait()

	l.connWG.Wait()
	newconn.Close()

	l.Close()
	listenerWG.Wait()

	// t.Logf("\nconn:\n%+v\n\n", conn)
	// t.Logf("\nconn.ConnectionState():\n%+v\n\n", conn.ConnectionState())
	// t.Logf("\ntls.Config:\n%+v\n\n", tlsConf)
}

// func TestClientDisconnect(t *testing.T) {
// 	// todo: we need to verify that events occur in the order that we want them (either via event hooks or log analysis)
// 	// this seems like a listener test than a integration test from a client perspective
// 	s := getServer(t)
// 	c := getClientConnWithClientCert(t)
// 	if err := c.Close(); err != nil {
// 		t.Fatal("Failed to close the client connection:", err)
// 	}
// 	if err := s.Stop(); err != nil {
// 		t.Fatal("Failed to stop the server:", err)
// 	}
// 	wg.Wait()
// }

func send(t *testing.T, conn *Conn, msg string) {
	n, err := conn.Write([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}
	if n < 100 {
		t.Logf("Sent message to listener from client: %v (%v bytes)", msg, n)
	} else {
		t.Logf("Sent message to listener from client: ... (%v bytes)", n)
	}
}

// closeGraceful waits for all request then connection handler goroutines to return then closes the listener. This method is meant for testing.
func closeGraceful(l *Listener) error {
	// todo: more proper way is to do TCPConn.CloseRead()/reqWG.Wait()/TCPConn.CloseWrite()/listener.Close()
	// but that requires using net.TCPListener/TCPConn and then upgrading to TLS (which is also good when supporting UnixSocket)
	l.reqWG.Wait()
	l.connWG.Wait()
	return l.listener.Close()
}
