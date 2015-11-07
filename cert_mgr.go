package titan

import (
	"crypto/x509/pkix"
	"time"

	"github.com/neptulon/ca"
)

// CertMgr handles creation of client certificates.
type CertMgr struct {
	clientCACert, clientCAKey []byte
}

// NewCertMgr initializes a new certificate manager with given client CA certificate and private key.
func NewCertMgr(clientCACert, clientCAKey []byte) CertMgr {
	return CertMgr{clientCACert: clientCACert, clientCAKey: clientCAKey}
}

// GenClientCert generates a client certificate.
func (c *CertMgr) GenClientCert(userID string) (cert []byte, key []byte, err error) {
	return ca.GenClientCert(pkix.Name{
		CommonName: userID,
	}, time.Hour, 3072, c.clientCACert, c.clientCAKey)
}
