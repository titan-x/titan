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
	s := getServer(t)
	defer s.Stop()
	c := getClientConn(t)
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

func getClientConn(t *testing.T) *Conn {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	addr := "localhost:" + Conf.App.Port
	c, err := Dial(addr, caCertBytes)
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

	return s
}

const (
	// host = localhost, cn = localhost, org = devastator
	caCert = `-----BEGIN CERTIFICATE-----
MIIEXzCCArOgAwIBAgIRALLzuCUen8F+N9nX7/9mMS4wCwYJKoZIhvcNAQELMCkx
EzARBgNVBAoTCmRldmFzdGF0b3IxEjAQBgNVBAMTCWxvY2FsaG9zdDAgFw0xNTA0
MDQxNjMxMzFaGA8yMzA1MDEyNDE2MzEzMVowKTETMBEGA1UEChMKZGV2YXN0YXRv
cjESMBAGA1UEAxMJbG9jYWxob3N0MIIBuDANBgkqhkiG9w0BAQEFAAOCAaUAMIIB
oAKCAZcAte1gt4dgj3rtb9UUSdsPRNLbqnN5BObBAVX/jW2e8AMusUf5MDPz2zN/
occVGvOEnlePtTYv1zC5SNx4zcaf2psSnHxxd2Obz3SKhg14u2r68VJWB29UDeAP
Fv5W3uVqQKYGG61M7Ks9/kx9wS4kcxSbcp88yLGL48vGxzq3GV68IyGMk/NNTYRZ
/Y00uEJ7g9SMaO5/dxpYAy+XFlQGGzwPegYO3u3VmhEa2Wor0IBXCQ4LvaGH12pn
0qf8sfEMc3jPY90SmDIU/uhs5aNmPJS5RqY29oNw2FHwKd6de7rTQwG5tE/+e9rG
ftHarvw6RI9mdu08Fmfha2wAMG8nDLXQ+zJ5KhGdWyfGtF9C9i2totN5VZePBM2I
vcGQjZIdjW1etDoXtsqc2nvQrpjI5qX3r9DhqfO4mjwx+ussyOWhwEwo1updret7
FNOfvedsDDqsRBQoiyq0jyp1Ko0CdlvgnkQv2ysOjWQDebW1Ndq7SgIfktxO8nV3
S3y8PZqa1tP3RBxFC10FhNy9U6CH/NVLnqHh1QIDAQABo1gwVjAOBgNVHQ8BAf8E
BAMCAKQwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQF
MAMBAf8wFAYDVR0RBA0wC4IJbG9jYWxob3N0MAsGCSqGSIb3DQEBCwOCAZcAYtJx
lCymlr0Cl0DSxbtHHbHJY8vK79FB4MQFo4lv9Z/jOny07o3ZjiYvA0zI2oIavQsT
/dvphbtn7WJL/o/sEpXftEBNyUBzyTB4cXM1PsqYAf0id8OF3wgKF4SBHZt5s5hM
P98IANnXh++AXFNlgRmyyhGNaxgfwj5VMmaqcGoLvFX3LNRwAOHbv7ZSGgt74JqN
WnIoxsewrX7NU61hQXCqdOJGKdBPPU1K1ytECHCvxZXKtZKDnQ8xKcq7dDMt6Sgz
TSH7srLr90k4hjF+KStwjZTnmFCGX6i8AB+fz3fw6F1oYWPbi0WOEsX0hE8iPK+G
PeDcn6iyQoDrG1z+x6J4KVrDqrh0IEMSoO6bopYIQSbTIGi0EymYw7ijhkxnzGBF
mvBXOAMauTUTbkbfnDb6NQsPPg9ZMobVfoDiCGq48H5OMLurlOzTesxxVZbkjv2J
5L63PvBFIpXEPemyfmG5ClqKe2lru4a6u2OdCZvHLM3zKEMAwjULwOjDxhpBpED7
V6sKSN9XohrDzdPqnsxuO+JvTA==
-----END CERTIFICATE-----`

	// RSA 3248 bits
	caKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

	// host = client.localhost, cn = client.localhost, org = devastator
	clientCert = `-----BEGIN CERTIFICATE-----
MIIEYDCCArSgAwIBAgIRAJHlIWkKmUXzD/tbnO2UCHIwCwYJKoZIhvcNAQELMCkx
EzARBgNVBAoTCmRldmFzdGF0b3IxEjAQBgNVBAMTCWxvY2FsaG9zdDAgFw0xNTA0
MDQxNjMxMzNaGA8yMzA1MDEyNDE2MzEzM1owMDETMBEGA1UEChMKZGV2YXN0YXRv
cjEZMBcGA1UEAxMQY2xpZW50LmxvY2FsaG9zdDCCAbgwDQYJKoZIhvcNAQEBBQAD
ggGlADCCAaACggGXAO+UcefREtwxrvTl6C3bCyvuB3R10xTn37XtZjVvq6Lti055
04iUFhcx6xwpjGTMaS8lhUm4V6mlFYBbRndNimBPERM1/k6u8VfunGo4JD0qJHgo
N5oq0RgqKQ6MS/AzE/tfjTUGRYDIJfZNwuQlGyHu+O5gxpzVrpGnpMqFOgQb1HVc
yntLI/11K74khATxigaxrimYMljVtS4hIC2qGadfobN225NzKwt+m3i/wZXIl9L2
SCFIhFt6YdzVWN3/JcUKHNpWVGXMqG1KNALJo2cL2a3drBouGaKKUXBSQw1SswmY
uUgLEKhOai0qmcz1jPuft87IvcDQt14o114gbB+miZIMIhONfv4x4dynPG5+zzDe
d+XR/1WOuN530bz8byuytUv7l3QfOdBsuJT7YX8p0m38V4xLu72cIH0PqDRn2PVX
7OBCz0lERH3RHZStARqnV6O8r1G1Sv0JNK0Lu5sOO/WifQG84xcrDmwiOjRe4+b4
0ISxHNE5G4YTv3FXVDn+xwlIkEwRLgbbfz1Bvrgo0ZMexZUCAwEAAaNSMFAwDgYD
VR0PAQH/BAQDAgCgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAw
GwYDVR0RBBQwEoIQY2xpZW50LmxvY2FsaG9zdDALBgkqhkiG9w0BAQsDggGXAFuj
MHQEeWm7lgK+AkGD6ZTALQR0JguXa6in+e8hTdwU7wqiOW3uBN9wgxDhhiXn012T
35XngoRp9W6nrOteWeRfvlCzWEzKxfTRm4fYahHthxR8w/hC7Mz2jowwbckKy77n
1K0JCZt5pkzrb/zyzykeQm+Ts+P41xjzYctnUrFLPAQTXRvKKe560LkSh4KO9QVy
vGHkFj3gAmHOA9CgmzEGsAe/ot5svvPEnNLPSNLyffzL37XtSUSzqvqrup+56c2B
ShR7dwWDNDqDqaavawrkiTegJ0ZfwZO5Qgk1asXpqvfekqtA49cmTSfu0W2JrjSb
2XU5MHW/UHXdeAiWsFapt3PNZBV6C/E/5m4yCLLaX2EeqIzUxlrDj+Ub6qlZ6qVM
LsyP0hprZj7b6G32Dkxtga8F9myvFDRExqvP8R8Z7c2JgVxRiPQ3XUrUFbVZn4+j
dGsVYK1VLkoOGTLTTqLEsDklljsx35H8uZuQKd34rngUMEYCOys6cMDr92rP6t5K
EMDfomB5c+U39lG5Ly9gTXCD2d8=
-----END CERTIFICATE-----`

	// RSA 3248 bits
	clientKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
)
