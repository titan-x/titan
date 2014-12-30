package main

import "testing"

func TestConfig(t *testing.T) {
	if Conf.App.Env != "test" || !Conf.App.Debug {
		t.Error("Config file is not initialized properly for development environment.")
	}

	// try the init block and see if it covers other files too regardless of execution order
}
