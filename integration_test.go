package main

import (
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
	s := getServer(t, false)
	c := getClientConn(t, true)
	if err := c.Close(); err != nil {
		t.Fatal("Failed to close the client connection:", err)
	}
	if err := s.Stop(); err != nil {
		t.Fatal("Failed to stop the server gradefully:", err)
	}

	// t.Fatal("Client method.close request was not handled properly")
	// t.Fatal("Client disconnect was not handled gracefully")
	// t.Fatal("Server method.close request was not handled properly")
	// t.Fatal("Server disconnect was not handled gracefully")
}

func getClientConn(t *testing.T, useClientCert bool) *Conn {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	var cert, key []byte
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

func getServer(t *testing.T, createNewCertPair bool) *Server {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	if createNewCertPair {
		var err error
		if caCertBytes, caKeyBytes, clientCertBytes, clientKeyBytes, err = genTestCertPair(512); err != nil {
			t.Fatal(err)
		}
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
MIIEXzCCArOgAwIBAgIRAMVGwitA8U2ZAKKq36q0Z9owCwYJKoZIhvcNAQELMCkx
EzARBgNVBAoTCmRldmFzdGF0b3IxEjAQBgNVBAMTCWxvY2FsaG9zdDAgFw0xNTA0
MDYwMDM1MzRaGA8yMzA1MDEyNjAwMzUzNFowKTETMBEGA1UEChMKZGV2YXN0YXRv
cjESMBAGA1UEAxMJbG9jYWxob3N0MIIBuDANBgkqhkiG9w0BAQEFAAOCAaUAMIIB
oAKCAZcAwxBVuSvSq4G/8MMhFujpvLSJWjv0O1m3xbskq0ZFMyHM3bttLf6ZFM/b
3LzFGbGiQddaFjtN+esmCYFrHrThPVjTz1R7YjtHl0Bky5v1Lhfm9ZngaVihsQLh
Ur4XsxneUHaHRQ+TKUeXm5n5bpNMQxlY1SfE4YJ4I1eFfpZTpKnaIGHzuj8oN1rb
eFqmXjWahngTVnyhUFL3PNpSIKtu7VcLGjt5xQVRtIhAbn6NHz6lzKawgN0qjZ0W
7A1g8Prdm1K2o7DETiJiF14kNyGoxOGf2PyUl0upyJjxKJxFbWEK6xlKglJILTM6
wFdvBbmR5gs2sKjN/ENAwRkX7/GGroNtOtda2vFVLrFC3T1XNpFpTVFyirBvRP+D
2dYGLKAbGdUL6DYYz0byvBzYOMyFpJ8mCmAmsVFzy+dSo67YGA1WyBTIfrsLFBpd
PgKEnPUovl6OcgWNwscxmKjpTX7K1vCp2LNUyo688nPQklkLHJq1PMtDE+KiBka8
ZHyYhD6HHVefPxp0AXlHUxflRM6hQWtMquHBXwIDAQABo1gwVjAOBgNVHQ8BAf8E
BAMCAKQwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQF
MAMBAf8wFAYDVR0RBA0wC4IJbG9jYWxob3N0MAsGCSqGSIb3DQEBCwOCAZcAuT9d
fN/YXLUsu0DB/7zEfjmpwQh1zKix/+mPJNuG6sdf9tJJSUfil0RwDlOZpnpJgIIH
PH/er/NfxcTmJNVofsc4ohrM3vVz1PfBobrEmEKaMjwPNMmtFIH1jIv+Ov8SVMZ7
BrZajam7HZPcoFLYBqco8YqjUp845+/cKz1AF0ND5uaZE7Q6rsNYMxYNEN/Bqabg
to7m/mF7h0TjSEMA5ScC+hhZ7VzwcdJhg3kRf7AaUjmgYCV3bAzeCIy/wQnu9iI2
kb3UK0BruRQK2zgAA8o5aY+JNyZxHsitTPL3dkV0VFUbpr3hxTnccydrRJoi9SVT
9E4WwH47PSKPw7/AKhGoH6aO+M1BvJj9h5O5gA82d38S4DPxvq+9hUZg2HhDHIax
1doFA5LCGU/5EDAdCZ3nyHSO5qla62gSPY2Yi6vG9S+TRTO3pJ14RdzTJv7NHorc
eWsi53SZpv/Sx6PgIHRxjWEaWXmi+bmZfjstDea5Cfwe1BlDQWELyU9G56fC7ksC
P0dAwe/mKAOsCX6DcjdUJAmZdg==
-----END CERTIFICATE-----`

	// RSA 3248 bits
	caKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

	// host = client.localhost, cn = client.localhost, org = devastator
	clientCert = `-----BEGIN CERTIFICATE-----
MIIEYDCCArSgAwIBAgIRAJ/2xbVwEwixc/0miP6OvOUwCwYJKoZIhvcNAQELMCkx
EzARBgNVBAoTCmRldmFzdGF0b3IxEjAQBgNVBAMTCWxvY2FsaG9zdDAgFw0xNTA0
MDYwMDM1NDRaGA8yMzA1MDEyNjAwMzU0NFowMDETMBEGA1UEChMKZGV2YXN0YXRv
cjEZMBcGA1UEAxMQY2xpZW50LmxvY2FsaG9zdDCCAbgwDQYJKoZIhvcNAQEBBQAD
ggGlADCCAaACggGXAOX6+uDatC8KMOm4hrV0W/2VLC4XEFVQOOjU3mSQcH3EX1Je
ym7yrOunxxv1v55wwdKyQK2Uvq8JsIyuyqn+NSL8bYXG64QQJbIt1QMLADZraX3E
EUTofbJ+7UAmdU899jLKKUZ/otdiJ/jjNy57COcOhyaxvDTTppqXR7eCd5AhpPyy
YN9l1E4IkmRi5IBeDMppSjHhveQrS2Nk+SrZCneJrfXbR4r8eLFHWZoYso2GTmqP
nkU9zok++JrXSew8HTnZAgLChdIb3mi6e3oOoNCYgwKeTDog/2wrqKU5XiDljuB/
KXjP6yLfxbh/NxvHEWNo99iIPl1Y0fzEdOhntIZaAiafd/sPKduxQwMYYLRcWngv
9eRauaaFPIOEUKWwQhPMjJwY9IFa8qqSjjWRMWpmoj3uhcK/V2csbA4iGj4PMdfa
yk+IWCSfNoJb8LGhCoNzNIDR9BpQP41XIqlsGVr4mprZGUjq18LvfMMhzF+08qeS
1gR3y2JFaSJqzGcgdmMfG4Dwx8j0+JWkRUpkbCLiwLrKLmMCAwEAAaNSMFAwDgYD
VR0PAQH/BAQDAgCgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAw
GwYDVR0RBBQwEoIQY2xpZW50LmxvY2FsaG9zdDALBgkqhkiG9w0BAQsDggGXAFIR
zoYGDBYVqpcBFICUzHyW6Nn10EYI0Soi7rKgtrXxu/Rhm7Odxol1sEq702CK533l
T/KkD9YFi2wJiy4HqbMxjs5bnkCyaHtRuTvQc3t+kicSlE+b7ugNhb9HY9Mv2oTI
p3FWzOpLid149VEwW13sirxRD8QFAK/th8WoPzaXr13+9SRvfzk7NuxNfFE3cZQS
0VNsWNAcVuTFa8ulfGmpSQiVrlEiXfkuAOy16iZUkVEBuXLSxb7aI/TbJAp5o6Pr
HhJnZuI6PC8GS3B8a8TM9vJ2YpIrfX6Gi4sJoaC4OuCTr4lMVzuCiRdSCfD8/yLC
j56YItjhkz8oz2V5Ym2um50gno90lICng9VeJSvUs0tcsD1SPCJZmBMVu2poLi+1
DdMCARLJWTmWkXx3iRFDSv1goVy+ltqD8sONCAeLZbA6lOlxzHnlo8w/SmTbA1Xg
guz+gL5MiV48/9fUkrqcV0beS1TtVo7Qntm0KhkoCoA9pufhbEM3TgP3AsOyHDI7
Y2x9XN10TL8Fqbq43Pq40XAlL3M=
-----END CERTIFICATE-----`

	// RSA 3248 bits
	clientKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
)
