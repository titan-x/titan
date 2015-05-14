package devastator

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
	// keyLength := 0 // used for internal test cert generation
	keyLength := 512

	caCert, caKey, clientCert, clientKey, err := genTestCertPair(keyLength)
	if err != nil {
		t.Fatal(err)
	}

	if keyLength == 0 {
		fmt.Println("CA cert:")
		fmt.Println(string(caCert))
		fmt.Println(string(caKey))
		fmt.Println("Client cert:")
		fmt.Println(string(clientCert))
		fmt.Println(string(clientKey))
	}
}

func genTestCertPair(keyLength int) (caCert, caKey, clientCert, clientKey []byte, err error) {
	// CA certificate
	caCert, caKey, err = genCert("127.0.0.1", 0, nil, nil, keyLength, "127.0.0.1", "devastator")

	if err != nil {
		err = fmt.Errorf("Failed to generate CA certificate or key: %v", err)
		return
	}
	if caCert == nil || caKey == nil {
		err = fmt.Errorf("Generated empty CA certificate or key")
		return
	}

	tlsCert, err := tls.X509KeyPair(caCert, caKey)

	if err != nil {
		err = fmt.Errorf("Generated invalid CA certificate or key: %v", err)
		return
	}
	if &tlsCert == nil {
		err = fmt.Errorf("Generated invalid CA certificate or key")
		return
	}

	// client certificate
	pub, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		err = fmt.Errorf("Failed to parse x509 certificate of CA cert to sign client-cert: %v", err)
		return
	}

	clientCert, clientKey, err = genCert("client.127.0.0.1", 0, pub, tlsCert.PrivateKey.(*rsa.PrivateKey), keyLength, "client.127.0.0.1", "devastator")
	if err != nil {
		err = fmt.Errorf("Failed to generate client-certificate or key: %v", err)
		return
	}

	tlsCert2, err := tls.X509KeyPair(clientCert, clientKey)

	if err != nil {
		err = fmt.Errorf("Generated invalid client-certificate or key: %v", err)
		return
	}
	if &tlsCert2 == nil {
		err = fmt.Errorf("Generated invalid client-certificate or key")
		return
	}

	return
}
