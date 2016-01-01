package test

import (
	"crypto/x509/pkix"
	"sync"
	"testing"
	"time"

	"github.com/nb-titan/titan"
	"github.com/neptulon/ca"
)

// ServerHelper is a titan.Server wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
	SeedData SeedData // Populated only when SeedDB() method is called.

	// PEM encoded X.509 certificate and private key pairs
	RootCACert,
	RootCAKey,
	IntCACert,
	IntCAKey,
	ServerCert,
	ServerKey []byte

	testing  *testing.T
	server   *titan.Server
	serverWG sync.WaitGroup // server instance goroutine wait group
	db       titan.InMemDB
}

// NewServerHelper creates a new server helper object.
// Titan server instance is initialized and ready to accept connection after this function return.
func NewServerHelper(t *testing.T) *ServerHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	// generate TLS certs
	certChain, err := ca.GenCertChain("FooBar", "127.0.0.1", "127.0.0.1", time.Hour, 512)
	if err != nil {
		t.Fatal("Failed to create TLS certificate chain:", err)
	}

	laddr := "127.0.0.1:" + titan.Conf.App.Port
	s, err := titan.NewServer(certChain.ServerCert, certChain.ServerKey, certChain.IntCACert, certChain.IntCAKey, laddr, titan.Conf.App.Debug)
	if err != nil {
		t.Fatal("Failed to create server:", err)
	}

	db := titan.NewInMemDB()
	if err := s.SetDB(db); err != nil {
		t.Fatal("Failed to attach InMemDB to server instance:", err)
	}

	h := ServerHelper{
		db:      db,
		server:  s,
		testing: t,

		RootCACert: certChain.RootCACert,
		RootCAKey:  certChain.RootCAKey,
		IntCACert:  certChain.IntCACert,
		IntCAKey:   certChain.IntCAKey,
		ServerCert: certChain.ServerCert,
		ServerKey:  certChain.ServerKey,
	}

	return &h
}

// SeedData is the data user for seeding the database.
type SeedData struct {
	User1 titan.User
	User2 titan.User
}

// SeedDB populates the database with the seed data.
func (sh *ServerHelper) SeedDB() *ServerHelper {
	cc1, ck1, err := ca.GenClientCert(pkix.Name{
		Organization: []string{"FooBar"},
		CommonName:   "1",
	}, time.Hour, 512, sh.IntCACert, sh.IntCAKey)
	if err != nil {
		sh.testing.Fatal(err)
	}

	cc2, ck2, err := ca.GenClientCert(pkix.Name{
		Organization: []string{"FooBar"},
		CommonName:   "2",
	}, time.Hour, 512, sh.IntCACert, sh.IntCAKey)
	if err != nil {
		sh.testing.Fatal(err)
	}

	sd := SeedData{
		User1: titan.User{ID: "1", Cert: cc1, Key: ck1},
		User2: titan.User{ID: "2", Cert: cc2, Key: ck2},
	}

	sh.db.SaveUser(&sd.User1)
	sh.db.SaveUser(&sd.User2)

	sh.SeedData = sd

	return sh
}

// Start starts the server.
func (sh *ServerHelper) Start() *ServerHelper {
	sh.serverWG.Add(1)
	go func() {
		defer sh.serverWG.Done()
		sh.server.Start()
	}()

	time.Sleep(time.Millisecond) // give Start() enough time to initiate
	return sh
}

// GetClientHelper creates and returns a ClientHelper that is connected to this server instance.
func (sh *ServerHelper) GetClientHelper() *ClientHelper {
	return NewClientHelper(sh.testing, sh.IntCACert, "127.0.0.1:"+titan.Conf.App.Port)
}

// Stop stops the server instance.
func (sh *ServerHelper) Stop() {
	if err := sh.server.Stop(); err != nil {
		sh.testing.Fatal("Failed to stop the server:", err)
	}
	sh.serverWG.Wait()
}
