package main

import "testing"

func TestConfig(t *testing.T) {
	config := GetConfig()
	if config.App.Env == "" {
		t.Error("Config file is not initialized.")
	}
}
