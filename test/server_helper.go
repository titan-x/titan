package test

import (
	"sync"
	"testing"
	"time"

	"github.com/nbusy/ca"
	"github.com/nbusy/devastator"
)

// ServerHelper is a devastator.Server wrapper.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
	db      devastator.InMemDB
	server  *devastator.Server
	testing *testing.T

	listenerWG sync.WaitGroup // server listener goroutine wait group

	// PEM encoded X.509 certificate and private key pairs
	RootCACert,
	RootCAKey,
	IntCACert,
	IntCAKey,
	ServerCert,
	ServerKey []byte
}

// NewServerHelper creates a new server helper object.
func NewServerHelper(t *testing.T) *ServerHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	// generate TLS certs
	certChain, err := ca.GenCertChain("FooBar", "127.0.0.1", "127.0.0.1", time.Hour, 512)
	if err != nil {
		t.Fatal("Failed to create TLS certificate chain:", err)
	}

	laddr := "127.0.0.1:" + devastator.Conf.App.Port
	s, err := devastator.NewServer(certChain.ServerCert, certChain.ServerKey, certChain.IntCACert, certChain.IntCAKey, laddr, devastator.Conf.App.Debug)
	if err != nil {
		t.Fatal("Failed to create server:", err)
	}

	db := devastator.NewInMemDB()
	if err := s.UseDB(db); err != nil {
		t.Fatal("Failed to attach InMemDB to server instance:", err)
	}

	h := ServerHelper{
		db:      db,
		server:  s,
		testing: t,

		RootCACert: certChain.RootCACert,
	}

	h.listenerWG.Add(1)
	go func() {
		defer h.listenerWG.Done()
		s.Start()
	}()

	return &h
}

// SeedData is the data user for seeding the database.
type SeedData struct {
	User1 devastator.User
	User2 devastator.User
}

// SeedDB populates the database with the seed data.
func (s *ServerHelper) SeedDB() *ServerHelper {
	// if certChain.ClientCert, certChain.ClientKey, err = ca.GenClientCert(pkix.Name{
	// 	Organization: []string{"FooBar"},
	// 	CommonName:   "1",
	// }, time.Hour, 512, certChain.IntCACert, certChain.IntCAKey); err != nil {
	// 	t.Fatal(err)
	// }
	//
	// if client2Cert, client2Key, err = ca.GenClientCert(pkix.Name{
	// 	Organization: []string{"FooBar"},
	// 	CommonName:   "2",
	// }, time.Hour, 512, certChain.IntCACert, certChain.IntCAKey); err != nil {
	// 	t.Fatal(err)
	// }

	s.db.SaveUser(&devastator.User{ID: "1"})
	s.db.SaveUser(&devastator.User{ID: "2"})
	return s
}

// Stop stops a server instance.
func (s *ServerHelper) Stop() {
	if err := s.server.Stop(); err != nil {
		s.testing.Fatal("Failed to stop the server:", err)
	}
	s.listenerWG.Wait()
}
