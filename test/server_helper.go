package test

import (
	"sync"
	"testing"

	"github.com/nbusy/devastator"
)

// ServerHelper is a devastator.Server wrapper with built-in error logging for testing.
type ServerHelper struct {
	DB         devastator.InMemDB
	server     *devastator.Server
	testing    *testing.T
	listenerWG sync.WaitGroup // server listener goroutine wait group
}

// NewServerHelper creates a new devastator.Server wrapper which has built-in error logging for testing.
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

	h := ServerHelper{DB: db, server: s, testing: t}

	h.listenerWG.Add(1)
	go func() {
		defer h.listenerWG.Done()
		s.Start()
	}()

	return &h
}

// SeedDB populates the database with seed data for testing.
func (s *ServerHelper) SeedDB() *ServerHelper {
	s.DB.SaveUser(&devastator.User{ID: 1, Cert: certChain.ClientCert})
	s.DB.SaveUser(&devastator.User{ID: 2})
	return s
}

// Stop stops a server instance with error checking.
func (s *ServerHelper) Stop() {
	if err := s.server.Stop(); err != nil {
		s.testing.Fatal("Failed to stop the server:", err)
	}
	s.listenerWG.Wait()
}
