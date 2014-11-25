package main

import "testing"

func TestConfig(t *testing.T) {
	if Conf.App.Env == "" {
		t.Error("Config file is not initialized.")
	}
}
