package test

import (
	"sync"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/titan-x/titan"
)

// ServerHelper is a titan.Server wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
	SeedData SeedData // Populated only when SeedDB() method is called.

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

	laddr := "127.0.0.1:" + titan.Conf.App.Port
	s, err := titan.NewServer(laddr)
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
	now := time.Now().Unix()
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims["userid"] = 1
	t.Claims["created"] = now
	ts1, err := t.SignedString(titan.Conf.App.JWTPass)
	t2 := jwt.New(jwt.SigningMethodHS256)
	t2.Claims["userid"] = 2
	t2.Claims["created"] = now
	ts2, err := t2.SignedString(titan.Conf.App.JWTPass)
	if err != nil {
		sh.testing.Fatalf("server-helper: failed to sign JWT token: %v", err)
	}

	sd := SeedData{
		User1: titan.User{ID: "1", JWT: ts1},
		User2: titan.User{ID: "2", JWT: ts2},
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
	return NewClientHelper(sh.testing, "127.0.0.1:"+titan.Conf.App.Port)
}

// Stop stops the server instance.
func (sh *ServerHelper) Stop() {
	if err := sh.server.Stop(); err != nil {
		sh.testing.Fatal("Failed to stop the server:", err)
	}
	sh.serverWG.Wait()
}
