package test

import (
	"sync"
	"testing"

	"github.com/nbusy/devastator"
)

// ServerHelper is a devastator.Server wrapper.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
	db         devastator.InMemDB
	server     *devastator.Server
	testing    *testing.T
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

	if certChain.RootCACert == nil {
		createCertChain(t)
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

	h := ServerHelper{db: db, server: s, testing: t}

	h.listenerWG.Add(1)
	go func() {
		defer h.listenerWG.Done()
		s.Start()
	}()

	return &h
}

// SeedDB populates the database with:
// - 2 users with their client certificates
func (s *ServerHelper) SeedDB() *ServerHelper {
	s.db.SaveUser(&devastator.User{ID: "1", Cert: certChain.ClientCert, Key: certChain.ClientKey})
	s.db.SaveUser(&devastator.User{ID: "2", Cert: client2Cert, Key: client2Key})
	return s
}

// Stop stops a server instance.
func (s *ServerHelper) Stop() {
	if err := s.server.Stop(); err != nil {
		s.testing.Fatal("Failed to stop the server:", err)
	}
	s.listenerWG.Wait()
}
