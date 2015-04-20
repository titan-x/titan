package main

import (
	"strconv"
	"testing"
)

func TestLen(t *testing.T) {
	a, _ := strconv.Atoi("12344324")
	t.Log(a)
}

func TestListener(t *testing.T) {
	host := "localhost:" + Conf.App.Port
	cert, privKey, _ := genCert("localhost", 0, nil, nil, 512, "localhost", "devastator")
	listener, err := Listen(cert, privKey, host, Conf.App.Debug)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go listener.Accept(func(conn *Conn, session *Session, msg []byte) {
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

	send(t, conn, "ping")
	send(t, conn, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	send(t, conn, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	send(t, conn, "In sit amet lectus felis, at pellentesque turpis.")
	send(t, conn, "Nunc urna enim, cursus varius aliquet ac, imperdiet eget tellus.")
	send(t, conn, randString(45000))
	send(t, conn, "close")

	t.Logf("\nconn:\n%+v\n\n", conn)
	t.Logf("\nconn.ConnectionState():\n%+v\n\n", conn.ConnectionState())
}

func send(t *testing.T, conn *Conn, msg string) {
	err := conn.Write([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}
}
