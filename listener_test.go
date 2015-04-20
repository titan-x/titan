package main

import (
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

	conn, err := Dial(host, cert, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte("ping"))
	conn.Write([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit."))
	conn.Write([]byte("In sit amet lectus felis, at pellentesque turpis."))
	conn.Write([]byte("Nunc urna enim, cursus varius aliquet ac, imperdiet eget tellus."))
	// conn.Write([]byte(randString(45000)))
	conn.Write([]byte("close"))

	t.Logf("\nconn:\n%+v\n\n", conn)
	t.Logf("\nconn.ConnectionState():\n%+v\n\n", conn.ConnectionState())

	wg.Wait()
}
