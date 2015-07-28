package test

import (
	"sync"
	"testing"

	"github.com/nbusy/devastator"
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

	createCertChain(t)

	laddr := "127.0.0.1:" + devastator.Conf.App.Port
	s, err := devastator.NewServer(caCertBytes, caKeyBytes, clientCACertBytes, clientCAKeyBytes, laddr, devastator.Conf.App.Debug)
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
