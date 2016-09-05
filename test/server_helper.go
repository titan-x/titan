package test

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/titan-x/titan"
	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/data/aws"
	"github.com/titan-x/titan/data/inmem"
)

var awsFlag = flag.Bool("aws", false, "Run tests with AWS support.")

// ServerHelper is a titan.Server wrapper for testing.
// All the functions are wrapped with proper test runner error logging.
type ServerHelper struct {
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
