package test

import (
	"crypto/x509/pkix"
	"testing"
	"time"

	"github.com/nbusy/ca"
)

// certificates for testing
var caCert, caKey, intermediateCACert, intermediateCAKey, serverCert, serverKey, clientCert, clientKey []byte

// create an entire certificate chain for testing: CA/SigningCert/HostingCert/ClientCert
func createCertChain(t *testing.T) {
	var err error

	caCert, caKey, err = ca.CreateCACert(pkix.Name{
		Organization:       []string{"Devastator"},
		OrganizationalUnit: []string{"Devastator Certificate Authority"},
		CommonName:         "Devastator Root CA",
	}, time.Hour, 512)

	if caCert == nil || caKey == nil || err != nil {
		t.Fatal("Failed to created CA cert", err)
	}

	intermediateCACert, intermediateCAKey, err = ca.CreateSigningCert(pkix.Name{
		Organization:       []string{"Devastator"},
		OrganizationalUnit: []string{"Devastator Intermediate Certificate Authority"},
		CommonName:         "Devastator Intermadiate CA",
	}, time.Hour, 512, caCert, caKey)

	if intermediateCACert == nil || intermediateCAKey == nil || err != nil {
		t.Fatal("Failed to created signing cert", err)
	}

	serverCert, serverKey, err = ca.CreateServerCert(pkix.Name{
		Organization: []string{"Devastator"},
		CommonName:   "127.0.0.1",
	}, "127.0.0.1", time.Hour, 512, intermediateCACert, intermediateCAKey)

	if serverCert == nil || serverKey == nil || err != nil {
		t.Fatal("Failed to created server cert", err)
	}

	clientCert, clientKey, err = ca.CreateClientCert(pkix.Name{
		Organization: []string{"Devastator"},
		CommonName:   "chuck.norris",
	}, time.Hour, 512, intermediateCACert, intermediateCAKey)

	if clientCert == nil || clientKey == nil || err != nil {
		t.Fatal("Failed to created client cert", err)
	}
}
