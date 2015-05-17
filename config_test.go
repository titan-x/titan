package devastator

import "testing"

func init() {
	Conf.App.Env = test
	Conf.App.Debug = true
}

func TestConfig(t *testing.T) {
	if Conf.App.Env != "test" || !Conf.App.Debug {
		t.Fatal("Config file is not initialized properly for testing environment")
	}
}

func TestGoEnv(t *testing.T) {
	// GO_ENV should be used if "DEVASTATOR_ENV" is empty
	// if both are empty, should be env = dev
}
