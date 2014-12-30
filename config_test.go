package main

import "testing"

func TestConfig(t *testing.T) {
	if Conf.App.Env != "test" {
		t.Error("Config file is not initialized properly for development environment.")
	}
}
