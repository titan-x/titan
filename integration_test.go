package main

import "testing"

func TestGoogleAuth(t *testing.T) {
	// t.Fatal("Google+ first sign-in (registration) failed with valid credentials")
	// t.Fatal("Google+ second sign-in (regular) failed with valid credentials")
	// t.Fatal("Google+ sign-in passed with invalid credentials")
	// t.Fatal("Authentication was not ACKed")
}

func TestClientCertAuth(t *testing.T) {
	// c := getConn(t)
	// t.Fatal("Authentication failed with a valid client certificate")
	// t.Fatal("Authenticated with invalid/expired client certificate")
	// t.Fatal("Authentication was not ACKed")
}

func TestTokenAuth(t *testing.T) {
	// t.Fatal("Authentication failed with a valid token")
	// t.Fatal("Authenticated with invalid/expired token")
	// t.Fatal("Authentication was not ACKed")
}

func TestReceiveOfflineQueue(t *testing.T) {
	// t.Fatal("Failed to receive queued messages after coming online")
	// t.Fatal("Failed to send ACK for received message queue")
}

func TestSendEcho(t *testing.T) {
	// t.Fatal("Failed to send a message to the echo user")
	// t.Fatal("Failed to send batch message to the echo user")
	// t.Fatal("Failed to send large message to the echo user")
	// t.Fatal("Did not receive ACK for a message sent")
	// t.Fatal("Failed to receive a response from echo user")
	// t.Fatal("Could not send an ACK for an incoming message")
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

func getConn(t *testing.T) *Conn {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	host := "localhost:" + Conf.App.Port
	c, err := Dial(host, []byte(cert))
	if err != nil {
		t.Fatal("Cannot connect to server:", err)
	}

	return c
}

const (
	// host = localhost, cn = localhost, org = devastator
	cert = `-----BEGIN CERTIFICATE-----
	MIIEXjCCArKgAwIBAgIQFvlJB+ijEORJb3nMgkF8NzALBgkqhkiG9w0BAQswKTET
	MBEGA1UEChMKZGV2YXN0YXRvcjESMBAGA1UEAxMJbG9jYWxob3N0MCAXDTE1MDQw
	NDE1MzA0NFoYDzIzMDUwMTI0MTUzMDQ0WjApMRMwEQYDVQQKEwpkZXZhc3RhdG9y
	MRIwEAYDVQQDEwlsb2NhbGhvc3QwggG4MA0GCSqGSIb3DQEBAQUAA4IBpQAwggGg
	AoIBlwC61PZSmiqwfgb/8z5q0bqoynSPtxIBMhpwcatrcUZXrIh2Y07v6inJgA+c
	8FQMScPeS6QtVSsfiv/90CLYzQB6+96sfCkn1fOzYr7kS2SlCAaH4WWL5gmSx0Fi
	vABtxFfm/GieCnaQC9JVnGedQCIOQTkZM3NaCNZzHU5GpQ+9rwN8oyFK/+Q1N5aH
	DANLeZnAMa869wIXt6CLjRZEH87FUThhdtb8/9NU/z4yGV+/IleS7XpmvokS/SBH
	3y8gogYPtgnWRfHE8qicuvb188uF+O118VtB34LHFKm5LMKMz46E1pc4DcGGQAHl
	LeCMDSaxZF8qLlYzAegAcx3jUM7fMquALmNwyT51Qf3HxhrCaSNxpdwnNZm0NEnQ
	KR3FfLTONDS25Ph2yAz/AfS2nhn2Ze/Llz3LMWlcNmFr5xOxmjxRU84ZC0VptTb7
	1d4YBoKGTO++AFN0ydR5+1FqRrmxsl2qAZAnNTnTsoX0SsE8JGVTwyWU90uGiCDo
	7LcRjSu9zZTK+GiJyWKh1XHO6iTZ8tb3JhodAgMBAAGjWDBWMA4GA1UdDwEB/wQE
	AwIApDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUw
	AwEB/zAUBgNVHREEDTALgglsb2NhbGhvc3QwCwYJKoZIhvcNAQELA4IBlwAtB6Zk
	B4CgPcz6tBbqzJzISM874WJq+gBP7Ju1eEbU4DU0md0XFXoHu1/62hJXtTOzqgOC
	RDBBqOLnmab78vS/LQL7o7MXfuOxQxSlvueBirclpk79MK2kxSpLC85I52XwgNNz
	jt2k1MztQOnx+ziO9CIzCcEWwdV5Rn+IqS95xbSui3geSPhGw3v5AmB5l4xLcD0o
	YZoKEz//J+LvucFbqoG5Avlq/EHigkf2xnStY596m9NK6KYwydoNyFzVk0FAeJET
	hUkDWBPaOFT8vlvAs28F0TYxWGcrDopV98GXmDHTW7sWOPJ1zyjrl9mtSD7TcS4R
	agbLy4ZiSmzqDIGWOah8cifwNTqUI4ROxGhkYgdi5vK6k3VKHI0d4y4sdkZiCHHx
	e6FgjzkXJtiDJu7J3e6SYS40MNKejcyOQw+C2+L+zRkJeBO7AJ3RRtwcb/UHwSbX
	YYty/fCBuifCXRqneZG43BvVvQyqXUprvIDxt1JQA99lg7D41icioq/G5meAvowF
	K+CBGFP2vWJ5lqPtxRID3jTV
	-----END CERTIFICATE-----`

	// RSA 3248 bits
	key = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
)
