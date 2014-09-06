package main

import "testing"

func TestConfig(t *testing.T) {
	config := GetConfig()
	if (config.App.Env != "development") {
		t.Error("Config file is not initialized.")
	}
}
