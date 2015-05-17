package devastator_test

// todo: should package be different to make this into a true integration test from a client perspective?

// this is the integration test package from a real client perspective.

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/nbusy/devastator"
)

var (
	caCertBytes     = []byte(caCert)
	caKeyBytes      = []byte(caKey)
	clientCertBytes = []byte(clientCert)
	clientKeyBytes  = []byte(clientKey)
)

func TestClientDisconnect(t *testing.T) {
	// todo: we need to verify that events occur in the order that we want them (either via event hooks or log analysis)
	// this seems like a listener test than a integration test from a client perspective
	s := getServer(t)
	c := getClientConnWithClientCert(t)
	if err := c.Close(); err != nil {
		t.Fatal("Failed to close the client connection:", err)
	}
	if err := s.Stop(); err != nil {
		t.Fatal("Failed to stop the server:", err)
	}
}

func TestClientClose(t *testing.T) {
	// t.Fatal("Client method.close request was not handled properly")
}

func TestSendClose(t *testing.T) {
	// t.Fatal("Server method.close request was not handled properly")
}

func TestServerDisconnect(t *testing.T) {
	// t.Fatal("Server disconnect was not handled gracefully")
}

func TestServerClose(t *testing.T) {
	s := getServer(t)
	c := getClientConnWithClientCert(t)
	if err := s.Stop(); err != nil {
		t.Fatal("Failed to stop the server:", err)
	}
	if err := c.Close(); err != nil {
		t.Fatal("Failed to close the client connection:", err)
	}

	// test what happens when there are outstanding connections and/or requests that are being handled
	// destroying queues and other stuff during Close() might cause existing request handles to malfunction
}

func TestAuth(t *testing.T) {
	// t.Fatal("Unauthorized clients cannot call any function other than method.auth.")
}

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
	// t.Fatal("Failed to stop the server: not all the goroutines were terminated properly")
	// t.Fatal("Failed to stop the server: server did not wait for ongoing read/write operations")
}

func TestConnTimeout(t *testing.T) {
	// t.Fatal("Send timout did not occur")
	// t.Fatal("Wait timeout did not occur")
	// t.Fatal("Read timeout did not occur")
}

func TestPing(t *testing.T) {
	// t.Fatal("Pong/ACK was not sent for ping")
}

func getClientConnWithClientCert(t *testing.T) *devastator.Conn {
	return _getClientConn(t, true)
}

func getAnonymousClientConn(t *testing.T) *devastator.Conn {
	return _getClientConn(t, false)
}

func _getClientConn(t *testing.T, useClientCert bool) *devastator.Conn {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	var cert, key []byte
	if useClientCert {
		cert = clientCertBytes
		key = clientKeyBytes
	}

	addr := "127.0.0.1:" + devastator.Conf.App.Port

	// retry connect in case we're operating on a very slow machine
	for i := 0; i <= 5; i++ {
		c, err := devastator.Dial(addr, caCertBytes, cert, key)
		if err != nil {
			if operr, ok := err.(*net.OpError); ok && operr.Op == "dial" && operr.Err.Error() == "connection refused" && i != 5 {
				time.Sleep(time.Millisecond * 50)
				continue
			}
			t.Fatalf("Cannot connect to server address %v with error: %v", addr, err)
		}

		if i != 0 {
			t.Logf("WARNING: it took %v retries to connect to the server, which might indicate code issues or slow machine.", i)
		}
		return c
	}
	panic("unreachable")
}

func getServer(t *testing.T) *devastator.Server {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	laddr := "127.0.0.1:" + devastator.Conf.App.Port
	s, err := devastator.NewServer(caCertBytes, caKeyBytes, laddr, devastator.Conf.App.Debug)
	if err != nil {
		t.Fatal("Failed to create server:", err)
	}

	go s.Start()
	return s
}

const (
	// host = 127.0.0.1, cn = 127.0.0.1, org = devastator
	caCert = `-----BEGIN CERTIFICATE-----
MIIEWjCCAq6gAwIBAgIRAPp9RGc+ttG6PwTzXx7JwWEwCwYJKoZIhvcNAQELMCkx
EzARBgNVBAoTCmRldmFzdGF0b3IxEjAQBgNVBAMTCTEyNy4wLjAuMTAgFw0xNTA1
MDkwOTA5NTZaGA8yMzA1MDIyODA5MDk1NlowKTETMBEGA1UEChMKZGV2YXN0YXRv
cjESMBAGA1UEAxMJMTI3LjAuMC4xMIIBuDANBgkqhkiG9w0BAQEFAAOCAaUAMIIB
oAKCAZcA645G4e2FK/69WW9nOEP9ZuXviitsa+MPhyW3N5JbMAWorMiRJwt0i/ZO
R8/nqiwe1Mk+OgavcEELygHi4MvPmwATO4E9uLiKBypmVeVA810a/8/1VJW2YZrv
N7MEcFWV2ItDRbsW9zCuAi9BznRPeHmEv5fC6+GnHHFj7ZU71hKJCB4jgR2xrrP+
5v3qKwe8bTZ0ZgcaLI0pSBWqgauLTvY6WOb8X8K1I2jmZ7SeH8QuMDJ8aCZ+Fz4P
S3QbUct95zlepftUDPuvgXf9Hrx+jMg1UADN3FL8NnB4TyXkBGQJd4W/D/QU/ykx
F/bLtis4v4gGla6CPxlfB8wI20XfOoJ+Q9AuCCn36xx94iZZNXsk0H8nx7tqa8L1
YD/DTzM7oGeOWINRUbMAlRyf8H/fWAN+kx6PXlafg60b83Ur+MTIC4sSoJRW+zEF
s7Om5I5tcZRdL92RHpbh+AIoTkwh0GIizTzM1b5R6Xfpxd2t6HCaK1+C3rz8LNvK
av9lrXpGQAitswzubZ0KxSXlsstFZluRAOxvKQIDAQABo1MwUTAOBgNVHQ8BAf8E
BAMCAKQwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQF
MAMBAf8wDwYDVR0RBAgwBocEfwAAATALBgkqhkiG9w0BAQsDggGXAJg35g6tv4wo
RrmFGeqCiv3j9LV8nMAG5KiSN3glUN3JSvdRAC5FU69iibeLUzNaroe7+ys9GM/7
CG4wwxmsYpu0tvEGw5nEN8sy5HmGuQv71yixnwHUBrvrEfGBjySL7M+NiCtCEKYg
qmsghfIXxenV9lcJdOlyoT84nXZq2kRaU+o+XO1tcH3nje5Ca0gC/SvIUNktB4mx
aRUhy9o8mWxSR9qy0u9c7gE4pLgC25bR7iI2VL01RYka1R049CVanChZM5jKSw71
XbdgkT1E4kM4JDUMJ6MTvvOrzCRPjujyOmbc6sIJG3xxvlnHqLz6j921rSlyGaXp
IFhoiMErWZPdttb+Z04ftdUd5TFvH5bXr49YBEQr2ClY4llfBDNYJC/YZs9evsed
iBxrgeMNw7M2VkAGS06zT0BxiIy3ez224BUgMq9QmFth7iah2o387y1lz2OFJIrq
Xfa/h41YVK3RPLXETVdTDzxUkaQRubmnTUaKzEfGUk5OB8lbbcFw1X1vmxPRKYjf
ED/9rwQ6WLC57vvbCoY=
-----END CERTIFICATE-----`

	// RSA 3248 bits
	caKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

	// host = client.127.0.0.1, cn = client.127.0.0.1, org = devastator
	clientCert = `-----BEGIN CERTIFICATE-----
MIIEXzCCArOgAwIBAgIQSwLu5wcVkGlY9qOW1pY4KTALBgkqhkiG9w0BAQswKTET
MBEGA1UEChMKZGV2YXN0YXRvcjESMBAGA1UEAxMJMTI3LjAuMC4xMCAXDTE1MDUw
OTA5MDk1OFoYDzIzMDUwMjI4MDkwOTU4WjAwMRMwEQYDVQQKEwpkZXZhc3RhdG9y
MRkwFwYDVQQDExBjbGllbnQuMTI3LjAuMC4xMIIBuDANBgkqhkiG9w0BAQEFAAOC
AaUAMIIBoAKCAZcAzTlkFCni6zUZza1X8UPkRLwlcMDBKpbxdQZrSQAhHqKZe7xg
ZthomSJ6Ahd0ueiXkGwDorhZLxW+1/iTQww4yDSMWmCNeflzuPp8E3maNTucYSSz
QGR4+GJx7+336spFBaT/ikGLHnbVaW6lGUAvKbRtUHFSoyfYit3Ar6x0+OnQrq+a
x4GFe+XiOvlZXqgKGQm1OWe58SFCbnvz+r0vbWIPabXk66gJMQ4yA0mjPKo7hEuE
h0XpUn2QiaSehS+NVmxjuM3j5fjPjMI4J3CXK4Ax9tM/imA2TbPig69T24CZpkHq
J/T99fgvEyO5+RcrrfHOjLHdVnyDceRXAB3avSOX8PiB3vF5kNruWN5GTLXxWaxW
IQkCtRjTWcN1APqFmkB0WEQcI+qgoE9WPSKUD2SMKzujd/HVn3xhKM3zAoSOZgBf
iXzQy6BRTxr7dE8/C/pultmVc4xS8YNGc0aIRGnht9s8yUSw7RQLb00Us4HE0RtX
FtRDEG+65V2d7yWtqN/ZQPo8832PfPSvR4XyrRUNp+Df9wIDAQABo1IwUDAOBgNV
HQ8BAf8EBAMCAKAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAb
BgNVHREEFDASghBjbGllbnQuMTI3LjAuMC4xMAsGCSqGSIb3DQEBCwOCAZcANbLw
HOQ2U84YnA30WUMhh6ETSJlwqOzSn2pyffI/EzjbBHYaQEWYWrZ7srNfF3+GctXj
rxyPodTJgzCrNFVAyE8V/Xm3DmCsxhJEGAapO/POFJ3wQNdYVK+yEox+lJDllz2Y
iffeV+WEV/6jixRsNDz5EXwljrUZiIeEXCWe+vpienOxB+Z+7Rh4JUbD+LuJCilx
XBs7uHSY8f2kCmu7lbI+5OrqO8lDhzGltBdt1cIB0T5dERaxn1JwWPuhCtyq2VQ+
1aaVD3t2E4ItlY4KNBW7tq9BLaO25cqwhIArWbuFY5ahV/hY/oJaNJI2V2D7f2YL
z6GHGVTA85rbYp9SLFKTQ9uT7OOL7LVburfKIgllgvgBPe1v+RwND7W6SWeSVZ/K
8CAQB4B8dzVCmpLbdvq0NBp0+8I2BZ0r+T9E42ddAnQqIt3x9dMj5jNX83LAuuzo
vrd2owpJtyFDW9U2uZrGqbzzRNTSAbcW5Oyf9PiCAxT8oDQ6YapEXr9dGRpMSGu/
ExyUsaeG4ZAGeYJp5v8Xda2CTw==
-----END CERTIFICATE-----`

	// RSA 3248 bits
	clientKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
)
