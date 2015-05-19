// Package test is an integration testing package for testing the server from a mobile client's perspective.
package test

import (
	"os"
	"testing"

	"github.com/nbusy/devastator"
)

// TestMain is the top level test runner with top level setup and teardown functions.
// These setup and teardown instructions are executed exactly once before and after all tests are run.
// http://golang.org/pkg/testing/#hdr-Main
func TestMain(m *testing.M) {
	// setup
	devastator.InitConf("test")
	// execute tests
	res := m.Run()
	// teardown
	// ...
	os.Exit(res)
}
