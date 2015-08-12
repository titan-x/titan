package test

import (
	"testing"
	"time"

	"github.com/nbusy/ca"
)

// certificates for testing
var certChain ca.CertChain

// create an entire certificate chain for testing: CA/SigningCert/HostingCert/ClientCert
func createCertChain(t *testing.T) {
	var err error
	if certChain, err = ca.GenCertChain("FooBar", "127.0.0.1", "127.0.0.1", time.Hour, 512); err != nil {
		t.Fatal(err)
	}
}
