package test

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/titan-x/titan"
	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/data/aws"
	"github.com/titan-x/titan/data/inmem"
	"github.com/titan-x/titan/models"
)

var awsFlag = flag.Bool("aws", false, "Run tests with AWS support.")

// ServerHelper is a titan.Server wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
	SeedData SeedData // Populated only when SeedDB() method is called.

	testing      *testing.T
	server       *titan.Server
	serverClosed chan bool
	db           data.DB
}

// NewServerHelper creates a new server helper object.
// Titan server instance is initialized and ready to accept connection after this function return.
func NewServerHelper(t *testing.T) *ServerHelper {
	if testing.Short() {
		t.Skip("Skipping integration test in short testing mode")
	}

	if (titan.Conf == titan.Config{}) {
		titan.InitConf("test")
	}

	url := "127.0.0.1:" + titan.Conf.App.Port
	s, err := titan.NewServer(url)
	if err != nil {
		t.Fatal("Failed to create server:", err)
	}

	var db data.DB
	if *awsFlag {
		db = aws.NewDynamoDB("", "")
	} else {
		db = inmem.NewDB()
	}
	if err := s.SetDB(db); err != nil {
		t.Fatal("Failed to attach InMemDB to server instance:", err)
	}

	h := ServerHelper{
		db:           db,
		server:       s,
		testing:      t,
		serverClosed: make(chan bool),
	}

	return &h
}

// SeedData is the data user for seeding the database.
type SeedData struct {
	User1 models.User
	User2 models.User
}

// SeedDB populates the database with the seed data.
func (sh *ServerHelper) SeedDB() *ServerHelper {
	now := time.Now().Unix()
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims["userid"] = "1"
	t.Claims["created"] = now
	ts1, err := t.SignedString([]byte(titan.Conf.App.JWTPass()))
	t2 := jwt.New(jwt.SigningMethodHS256)
	t2.Claims["userid"] = "2"
	t2.Claims["created"] = now
	ts2, err := t2.SignedString([]byte(titan.Conf.App.JWTPass()))
	if err != nil {
		sh.testing.Fatalf("server-helper: failed to sign JWT token: %v", err)
	}

	sd := SeedData{
		User1: models.User{ID: "1", JWTToken: ts1},
		User2: models.User{ID: "2", JWTToken: ts2},
	}

	sh.db.SaveUser(&sd.User1)
	sh.db.SaveUser(&sd.User2)

	sh.SeedData = sd

	return sh
}

// ListenAndServe starts the server.
func (sh *ServerHelper) ListenAndServe() *ServerHelper {
	go func() {
		if err := sh.server.ListenAndServe(); err != nil {
			sh.testing.Fatalf("server-helper: failed to start server: %v", err)
		}
		sh.serverClosed <- true
	}()

	time.Sleep(time.Millisecond) // give Start() enough time to initiate
	return sh
}

// GetClientHelper creates and returns a ClientHelper that is connected to this server instance.
func (sh *ServerHelper) GetClientHelper() *ClientHelper {
	return NewClientHelper(sh.testing, "ws://127.0.0.1:"+titan.Conf.App.Port)
}

// CloseWait closes the server and wait for all request/conn goroutines to exit.
func (sh *ServerHelper) CloseWait() {
	if err := sh.server.Close(); err != nil {
		sh.testing.Fatal("Failed to stop the server:", err)
	}
	select {
	case <-sh.serverClosed:
	case <-time.After(time.Second):
		sh.testing.Fatal("server didn't close in time")
	}

	if os.Getenv("TRAVIS") != "" || os.Getenv("CI") != "" {
		time.Sleep(time.Millisecond * 50)
	} else {
		time.Sleep(time.Millisecond * 5)
	}
}
