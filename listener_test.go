package main

import (
	"crypto/tls"
	"io"
	"strconv"
	"sync"
	"testing"
)

func TestLen(t *testing.T) {
	a, _ := strconv.Atoi("12344324")
	t.Log(a)
}

func TestListener(t *testing.T) {
	var wg sync.WaitGroup
	host := "localhost:" + Conf.App.Port
	cert, privKey, _ := genCert("localhost", 0, nil, nil, 512, "localhost", "devastator")
	listener, err := Listen(cert, privKey, host, Conf.App.Debug)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go listener.Accept(func(conn *Conn, session *Session, msg []byte) {
		wg.Add(1)
		defer wg.Done()
		// todo: compare sent/incoming messages for equality

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) > 0 {
			t.Logf("Client connected with client certificate subject: %v\n", certs[0].Subject)
		}
	}, func(conn *Conn, session *Session) {
	})

	conn, err := Dial(host, cert)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	send(t, conn, "4\nping")
	send(t, conn, "56\nLorem ipsum dolor sit amet, consectetur adipiscing elit.")
	send(t, conn, "49\nIn sit amet lectus felis, at pellentesque turpis.")
	send(t, conn, "64\nNunc urna enim, cursus varius aliquet ac, imperdiet eget tellus.")
	send(t, conn, "45000\n"+randString(45000))
	send(t, conn, "5\nclose")

	// t.Logf("\nconn:\n%+v\n\n", conn)
	// t.Logf("\nconn.ConnectionState():\n%+v\n\n", conn.ConnectionState())
	// t.Logf("\ntls.Config:\n%+v\n\n", tlsConf)

	wg.Wait()
}

func send(t *testing.T, conn *tls.Conn, msg string) {
	n, err := io.WriteString(conn, msg)
	if err != nil {
		t.Fatalf("Error while writing message to connection %v", err)
	}
	t.Logf("Sending message to listener from client: %v (%v bytes)", msg, n)
}
