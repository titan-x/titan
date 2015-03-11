package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLen(t *testing.T) {
	a, _ := strconv.Atoi("12344324")
	t.Log(a)
}

func TestListener(t *testing.T) {
	var wg sync.WaitGroup
	host := "localhost:" + Conf.App.Port
	cert, privKey, _ := genCert()
	listener, err := Listen(cert, privKey, host, Conf.App.Debug)
	if err != nil {
		t.Fatal(err)
	}

	go listener.Accept(func(conn *tls.Conn, session *Session, msg []byte) {
		wg.Add(1)
		defer wg.Done()
		t.Logf("Incoming message to listener from a client: %v", string(msg))
	}, func(conn *tls.Conn, session *Session) {
	})

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		panic("failed to parse root certificate")
	}

	conn, err := tls.Dial("tcp", host, &tls.Config{RootCAs: roots})
	if err != nil {
		t.Fatal(err)
	}

	send(t, conn, "4\nping")
	send(t, conn, "56\nLorem ipsum dolor sit amet, consectetur adipiscing elit.")
	send(t, conn, "49\nIn sit amet lectus felis, at pellentesque turpis.")
	send(t, conn, "64\nNunc urna enim, cursus varius aliquet ac, imperdiet eget tellus.")
	send(t, conn, "45000\n"+randString(45000))
	send(t, conn, "5\nclose")

	wg.Wait()
	time.Sleep(1000 * time.Millisecond) // todo: a more proper wait..
	conn.Close()
	listener.Close()
}

func send(t *testing.T, conn *tls.Conn, msg string) {
	n, err := io.WriteString(conn, msg)
	if err != nil {
		t.Fatalf("Error while writing message to connection %v", err)
	}
	t.Logf("Sending message to listener from client: %v (%v bytes)", msg, n)
}
