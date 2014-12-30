package main

import "testing"

func init() {
	Conf.App.Env = test
	Conf.App.Debug = true
}

func TestConfig(t *testing.T) {
	if Conf.App.Env != "test" || !Conf.App.Debug {
		t.Error("Config file is not initialized properly for development environment.")
	}
}
