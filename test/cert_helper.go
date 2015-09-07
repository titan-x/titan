package test

import (
	"crypto/x509/pkix"
	"testing"
	"time"

	"github.com/nbusy/ca"
)

// certificates for testing
var certChain ca.CertChain
var client2Cert, client2Key []byte

// create an entire certificate chain for testing: CA/SigningCert/HostingCert/ClientCert
func createCertChain(t *testing.T) {
	var err error
	if certChain, err = ca.GenCertChain("FooBar", "127.0.0.1", "127.0.0.1", time.Hour, 512); err != nil {
		t.Fatal(err)
	}

	if certChain.ClientCert, certChain.ClientKey, err = ca.GenClientCert(pkix.Name{
		Organization: []string{"FooBar"},
		CommonName:   "1",
	}, time.Hour, 512, certChain.IntCACert, certChain.IntCAKey); err != nil {
		t.Fatal(err)
	}

	if client2Cert, client2Key, err = ca.GenClientCert(pkix.Name{
		Organization: []string{"FooBar"},
		CommonName:   "2",
	}, time.Hour, 512, certChain.IntCACert, certChain.IntCAKey); err != nil {
		t.Fatal(err)
	}
}
