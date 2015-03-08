package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	mrand "math/rand"
	"net"
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
	cert, privKey := genCert(t)
	listener, err := Listen(cert, privKey, "localhost:8091", true)
	if err != nil {
		t.Fatal(err)
	}

	go listener.Accept(func(msg []byte) {
		wg.Add(1)
		defer wg.Done()
		t.Logf("Incoming message to listener from a client: %v", string(msg))
	})

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		panic("failed to parse root certificate")
	}

	conn, err := tls.Dial("tcp", "localhost:8091", &tls.Config{RootCAs: roots})
	if err != nil {
		t.Fatal(err)
	}

	send(t, conn, "4   ping")
	send(t, conn, "56  Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	send(t, conn, "49  In sit amet lectus felis, at pellentesque turpis.")
	send(t, conn, "64  Nunc urna enim, cursus varius aliquet ac, imperdiet eget tellus.")
	// send(t, conn, "8671"+randSeq(8671))
	send(t, conn, "5   close")

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
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

// Generate a self-signed PEM encoded X.509 certificate and private key pair (i.e. 'cert.pem', 'key.pem').
// Based on the sample from http://golang.org/src/crypto/tls/generate_cert.go (taken at Jan 30, 2015).
func genCert(t *testing.T) (pemBytes, privBytes []byte) {
	hosts := []string{"localhost"}
	privKey, err := rsa.GenerateKey(rand.Reader, 512)
	notBefore := time.Now()
	notAfter := notBefore.Add(290 * 365 * 24 * time.Hour) //290 years
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		t.Fatalf("failed to generate serial number: %s", err)
	}

	cert := x509.Certificate{
		IsCA:         true,
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"localhost"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			cert.IPAddresses = append(cert.IPAddresses, ip)
		} else {
			cert.DNSNames = append(cert.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, &privKey.PublicKey, privKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %s", err)
	}

	pemBytes = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privBytes = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}
