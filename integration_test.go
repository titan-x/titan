package main

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"testing"
)

var (
	caCertBytes     = []byte(caCert)
	caKeyBytes      = []byte(caKey)
	clientCertBytes = []byte(clientCert)
	clientKeyBytes  = []byte(clientKey)
)

func TestGoogleAuth(t *testing.T) {
	// t.Fatal("Google+ first sign-in (registration) failed with valid credentials")
	// t.Fatal("Google+ second sign-in (regular) failed with valid credentials")
	// t.Fatal("Google+ sign-in passed with invalid credentials")
	// t.Fatal("Authentication was not ACKed")
}

func TestClientCertAuth(t *testing.T) {
	keyLength := 512
	pemBytes, privBytes, err := genCert("localhost", 0, nil, nil, keyLength, "localhost", "devastator")
	tlsCert, err := tls.X509KeyPair(pemBytes, privBytes)
	pub, err := x509.ParseCertificate(tlsCert.Certificate[0])
	pemBytes2, privBytes2, err := genCert("client.localhost", 0, pub, tlsCert.PrivateKey.(*rsa.PrivateKey), keyLength, "client.localhost", "devastator")
	if err != nil {
		t.Fatal(err)
	}

	caCertBytes = pemBytes
	caKeyBytes = privBytes
	clientCertBytes = pemBytes2
	clientKeyBytes = privBytes2

	s := getServer(t)
	defer s.Stop()
	c := getClientConn(t, true)
	defer c.Close()

	// t.Fatal("Authentication failed with a valid client certificate")
	// t.Fatal("Authenticated with invalid/expired client certificate")
	// t.Fatal("Authentication was not ACKed")
}

func TestTokenAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid token")
	// t.Fatal("Authenticated with invalid/expired token")
	// t.Fatal("Authentication was not ACKed")
}

func BenchmarkAuth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func BenchmarkClientCertAuth(b *testing.B) {
	// for various certificate key sizes (512....4096) and ECDSA, and with/without resumed handshake / session tickets
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func TestReceiveOfflineQueue(t *testing.T) {
	// t.Fatal("Failed to receive queued messages after coming online")
	// t.Fatal("Failed to send ACK for received message queue")
}

func BenchmarkQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func TestSendEcho(t *testing.T) {
	// t.Fatal("Failed to send a message to the echo user")
	// t.Fatal("Failed to send batch message to the echo user")
	// t.Fatal("Failed to send large message to the echo user")
	// t.Fatal("Did not receive ACK for a message sent")
	// t.Fatal("Failed to receive a response from echo user")
	// t.Fatal("Could not send an ACK for an incoming message")
}

func BenchmarkParallelThroughput(b *testing.B) {
	// for various conn levels vs. message per second: 50:xxxx, 500:xxx, 5000:xx, ... conn/mps (hopefully!)
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func TestStop(t *testing.T) {
	// todo: this should be a listener/queue test if we don't use any goroutines in the Server struct methods
	// t.Fatal("Failed to stop the server gracefully: not all the goroutines were terminated properly")
	// t.Fatal("Failed to stop the server gracefully: server did not wait for ongoing read/write operations")
}

func TestConnTimeout(t *testing.T) {
	// t.Fatal("Send timout did not occur")
	// t.Fatal("Wait timeout did not occur")
	// t.Fatal("Read timeout did not occur")
}

func TestPing(t *testing.T) {
	// t.Fatal("Pong/ACK was not sent for ping")
}

func TestDisconnect(t *testing.T) {
	// t.Fatal("Client method.close request was not handled properly")
	// t.Fatal("Client disconnect was not handled gracefully")
	// t.Fatal("Server method.close request was not handled properly")
	// t.Fatal("Server disconnect was not handled gracefully")
}

func getClientConn(t *testing.T, useClientCert bool) *Conn {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	var cert []byte
	var key []byte
	if useClientCert {
		cert = clientCertBytes
		key = clientKeyBytes
	}

	addr := "localhost:" + Conf.App.Port
	c, err := Dial(addr, caCertBytes, cert, key)
	if err != nil {
		t.Fatalf("Cannot connect to server address %v with error: %v", addr, err)
	}

	return c
}

func getServer(t *testing.T) *Server {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	laddr := "localhost:" + Conf.App.Port
	s, err := NewServer(caCertBytes, caKeyBytes, laddr, Conf.App.Debug)
	if err != nil {
		t.Fatal("Failed to create server", err)
	}

	go s.Start()
	return s
}

const (
	// host = localhost, cn = localhost, org = devastator
	caCert = `-----BEGIN CERTIFICATE-----
MIIEVDCCAqigAwIBAgIQNQ40gKV/XmiDUV3rlLil1zALBgkqhkiG9w0BAQswKTET
MBEGA1UEChMKZGV2YXN0YXRvcjESMBAGA1UEAxMJbG9jYWxob3N0MCAXDTE1MDQw
NjAwMDQ0MloYDzIzMDUwMTI2MDAwNDQyWjApMRMwEQYDVQQKEwpkZXZhc3RhdG9y
MRIwEAYDVQQDEwlsb2NhbGhvc3QwggG4MA0GCSqGSIb3DQEBAQUAA4IBpQAwggGg
AoIBlwClgKZqFVj8vev2pKOS1rj5QWgmcM+GHTEzfQLYoW6Edug7G4AMy/C3Yvtz
7ZHyANiSgsvMoviCLgNDJaSY61eE1KVj+mvwSqMDA0CUiuz2SKmiQIPUEoU73oYe
FmNtNHNssjIGAx6Yel7K4S87YOSv5TIy3HQT10qIaTHLVOUNlEK0G9kY/pm7MKIC
RP2Q+D+P9vqqFlCQ4L4gUxXC/smnxzilEQ+8oB7eMuf3ym61p/zUE0aSfJnqWErc
rNd09HU5O1C49PlvkH0e5ov8WYBI3q/yVR63F8x6RcamWJhPKb6fLLkkx49FupRa
kRKECZbFQIWwtXrPPMXfobojt4eusoyOA3hVKai8O7DskjPuBMVEsKxLirnXhGO0
XqiR444LmgYuh7oORLr+FY8zS3/XqzK25DPp+rYPAi1ssChvSfBENf14vhjUufx3
J73sdRws7/9awigQLJdKPPGRB1or6GWn3zqmTLe4ZV5syKNoyZP+TdhRBD7hdbRv
A8+FonzPZV1S+bRUNrhXTQ86yIHEo5r02BjPAgMBAAGjTjBMMA4GA1UdDwEB/wQE
AwIApDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MBQGA1Ud
EQQNMAuCCWxvY2FsaG9zdDALBgkqhkiG9w0BAQsDggGXAI4mXvZSVXANqJI/7C8D
M+tmRzX0nRhh/ebGHXhWySa+ociNqxSmV9n/HyOdzCWciNF3LM4lpSeApuoD3zt0
94YWXPCTFtq9LvxBCmSil28mL8SCR3XXG0FHgh0LBaDRRYlpNv+KYO2i0eBGh8C3
PddKEtvCVEr0wPSmPveg2jKM1vpcsKX31KB/ORC+LJ9l6vQ83u/Kyz6Y0C054oQ/
1ru3NUCQnSbpIPqrPwngkBcACoO9ilfCJa7/tusOoUsKcjpX+1xNqDMBocvBWytx
sjvwxgl6U+h4KZ8B70iReaQBhU0L3eEOD6xS9KfiXDLszvsY4hR9Ileaj2eB8Htb
RD7yPW/G0YI0aCjwN0oJibOqTUs+764Z2cPa8DecWuTXYP4Uc8ICaCC6oG9Dchok
AtE7mN+zthLNfEPnE4lPeqs7xzoPVDUYvWLZnzgqq1bicuCghQPlpMe01vWlG9rW
s+/Et/UZns+LM73Zh6fRupkFWRWHQxi3Lq6hg4Q7uwFrjYQqH1PMRkEidAOsBhVt
AoYzzwA+y/I=
-----END CERTIFICATE-----`

	// RSA 3248 bits
	caKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

	// host = client.localhost, cn = client.localhost, org = devastator
	clientCert = `-----BEGIN CERTIFICATE-----
MIIEYDCCArSgAwIBAgIRAKVgVcSNf5fiKlcOKXnCaacwCwYJKoZIhvcNAQELMCkx
EzARBgNVBAoTCmRldmFzdGF0b3IxEjAQBgNVBAMTCWxvY2FsaG9zdDAgFw0xNTA0
MDYwMDA0NDNaGA8yMzA1MDEyNjAwMDQ0M1owMDETMBEGA1UEChMKZGV2YXN0YXRv
cjEZMBcGA1UEAxMQY2xpZW50LmxvY2FsaG9zdDCCAbgwDQYJKoZIhvcNAQEBBQAD
ggGlADCCAaACggGXAKvwi5ptcr/nKh4K6w5t1rd07RpwAr9/rMmI+0Y1/NOCBUCZ
/2RvMACb/DzB7OYsTBG+wUIN7yNeiSTbg6VItCIVVadRG231tiHwkjAdH75kNdFN
4oLlBz3xnqL0Js34WO71Py2d+WbAjFz5DOwYS3Pgm48lfMOAaQU6penl8mR4/L8G
B7/OzWSHMD7747FosBMxeLU5DGmRzNJ41SJhxU/Hl5sTwYXvw6ijsLdxLLonAQU+
NQ2MtLTvO4or0I3y89bioBF4tggz0cDku9laQgTecJykWrs6l+/g9eZHw4HFj39H
iT0X2rbxPPzBcQj+p/nkDo86/KhsVOYyDxosT+Wa1ek9opxcKOTJbv5LmS7BkWgH
1dbZxPTuF+HeAmvrykyaJyih4IBH2FtFjUP0pQSveOOnCszecpJhs6LM8rviCMrF
7Zwsm9yCUzeQEmWG6jt3bkk9+RQmZ3lcX+C6ysahkxiTdzL4OcQpXsBoDJByEtxF
i4R4Rtvn0HApYwvljBBTG/q0nOqpF8yo0dj2AzjjFlPA5wcCAwEAAaNSMFAwDgYD
VR0PAQH/BAQDAgCgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAw
GwYDVR0RBBQwEoIQY2xpZW50LmxvY2FsaG9zdDALBgkqhkiG9w0BAQsDggGXAEh5
zhvN2RyotXs/IOyR5g9c1VNfPpNEEwQLaYajnDXV7P7AK2UEsaD/Y1XvwoLJoCy8
2DzOqKUJldq8rLK2Ht2v1ToUb00mZ3ziXFjeDwURz+lcYi5QXxX1M1+LNd5oNAX3
ORFZVSo7eKtrKkRwG23BGS5fwEalAHVbl//+lraff0E2e5Zboe1/q1Uvo7Aam1vC
RmwuSij3X8gJD2RgL8GrHV3DDpFnjjzTxy6yhd5WVPbT1i7hUrj8qImeoaBz9NOw
7YdrrKEyQt4Ruhrc5CuMPu4dPsQOdXuzDz1WK3h1uhFwC85t/ESnPx9gPDxPEaPr
M5c00p+r85P5xgboBI7PXD4va2+j/4AgZPjz+4dZpDNsEKNCJeRWOKlZJiCLi1f8
Jx5hnQIeaS9mjMLenoSiQUphI5SrXe8td6xAHlJ6hcsdC5EFzSJPVgcpVRSetkjy
oUN0Bh0cyekxPLcGsqyLmONQtONKo7biuXwT5QwMW4R0VTSMa735bSf2vbRf0DRL
uXP6kSVfabFJ+jxqwDfYnFyHI20=
-----END CERTIFICATE-----`

	// RSA 3248 bits
	clientKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
)
