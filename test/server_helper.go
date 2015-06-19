package test

import (
	"sync"
	"testing"

	"github.com/nbusy/devastator"
)

// server certificate for testing
var (
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

	caCertBytes = []byte(caCert)
	caKeyBytes  = []byte(caKey)
)

// ServerHelper is a devastator.Server wrapper with built-in error logging for testing.
type ServerHelper struct {
	server     *devastator.Server
	testing    *testing.T
	listenerWG sync.WaitGroup // server listener goroutine wait group
}

// NewServerHelper creates a new devastator.Server wrapper which has built-in error logging for testing.
func NewServerHelper(t *testing.T) *ServerHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	laddr := "127.0.0.1:" + devastator.Conf.App.Port
	s, err := devastator.NewServer(caCertBytes, caKeyBytes, laddr, devastator.Conf.App.Debug)
	if err != nil {
		t.Fatal("Failed to create server:", err)
	}

	h := ServerHelper{server: s, testing: t}

	h.listenerWG.Add(1)
	go func() {
		defer h.listenerWG.Done()
		s.Start()
	}()

	return &h
}

// Stop stops a server instance with error checking.
func (s *ServerHelper) Stop() {
	if err := s.server.Stop(); err != nil {
		s.testing.Fatal("Failed to stop the server:", err)
	}
	s.listenerWG.Wait()
}
