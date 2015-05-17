package devastator

import (
	"crypto/tls"
	"crypto/x509"
	"strconv"
	"testing"
	"time"
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
	msg4 := randString(45000)
	msg5 := randString(500000)

	host := "localhost:" + Conf.App.Port
	cert, privKey, _ := genCert("localhost", 0, nil, nil, 512, "localhost", "devastator")
	l, err := Listen(cert, privKey, host, false)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	go l.Accept(func(conn *Conn, session *Session, msg []byte) {
		certs := conn.ConnectionState().PeerCertificates
		if len(certs) > 0 {
			t.Logf("Client connected with client certificate subject: %v\n", certs[0].Subject)
		}

		m := string(msg)
		if m != msg1 && m != msg2 && m != msg3 && m != msg4 && m != msg5 {
			t.Fatal("Sent and incoming message did not match for message:", m)
		}
	}, func(conn *Conn, session *Session) {
	})

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
	defer l.connWG.Wait()
	defer conn.Close()

	newconn := NewConn(conn, 0, 0, 0, false)

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

	time.Sleep(time.Millisecond * 100)

	// t.Logf("\nconn:\n%+v\n\n", conn)
	// t.Logf("\nconn.ConnectionState():\n%+v\n\n", conn.ConnectionState())
	// t.Logf("\ntls.Config:\n%+v\n\n", tlsConf)
}

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
