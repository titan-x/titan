package main

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	l := 12304
	str := randString(l)

	if len(str) != l {
		t.Fatalf("Expected a random string of length %v but got %v", l, len(str))
	}
	if str[1] == str[2] && str[3] == str[4] && str[5] == str[6] && str[7] == str[8] {
		t.Fatal("Expected a random string, got repeated characters")
	}
}

func TestGetID(t *testing.T) {
	for i := 0; i < 50; i++ {
		id, err := getID()

		if err != nil {
			t.Fatalf("Error while generating unique ID: %v", err)
		}
		if len(id) != 26 {
			t.Fatalf("Expected a string of length 26 but got %v", len(id))
		}
		if id[3] == id[4] && id[5] == id[6] && id[7] == id[8] && id[9] == id[10] {
			t.Fatal("Expected a random string, got repeated characters")
		}
	}
}

func TestGenCert(t *testing.T) {
	keyLength := 0

	// CA certificate
	pemBytes, privBytes, err := genCert("localhost", 0, nil, nil, keyLength, "localhost", "devastator")

	if err != nil {
		t.Fatalf("Failed to generate CA certificate or key: %v", err)
	}
	if pemBytes == nil || privBytes == nil {
		t.Fatal("Generated empty CA certificate or key")
	}

	tlsCert, err := tls.X509KeyPair(pemBytes, privBytes)

	if err != nil {
		t.Fatalf("Generated invalid CA certificate or key: %v", err)
	}
	if &tlsCert == nil {
		t.Fatal("Generated invalid CA certificate or key")
	}

	// client certificate
	pub, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		t.Fatal("Failed to parse x509 certificate of CA cert to sign client-cert:", err)
	}

	pemBytes2, privBytes2, err := genCert("client.localhost", 0, pub, tlsCert.PrivateKey.(*rsa.PrivateKey), keyLength, "client.localhost", "devastator")
	if err != nil {
		t.Fatal("Failed to generate client-certificate or key:", err)
	}

	tlsCert2, err := tls.X509KeyPair(pemBytes2, privBytes2)

	if err != nil {
		t.Fatalf("Generated invalid client-certificate or key: %v", err)
	}
	if &tlsCert2 == nil {
		t.Fatal("Generated invalid client-certificate or key")
	}

	if keyLength == 0 {
		fmt.Println("CA cert:")
		fmt.Println(string(pemBytes))
		fmt.Println(string(privBytes))
		fmt.Println("Client cert:")
		fmt.Println(string(pemBytes2))
		fmt.Println(string(privBytes2))
	}
}
