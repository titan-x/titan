package titan

import "testing"

func init() {
	InitConf("test")
}

func TestConfig(t *testing.T) {
	// at this point, configuration must have been initialized with test environment defaults
	if Conf.App.Env != "test" || !Conf.App.Debug || Conf.App.Port != "3001" {
		t.Fatal("Config file is not initialized properly for testing environment")
	}

	InitConf("development")
	if Conf.App.Env != "development" || !Conf.App.Debug || Conf.App.Port != "3000" {
		t.Fatal("Config file is not initialized properly for development environment")
	}

	InitConf("production")
	if Conf.App.Env != "production" || Conf.App.Debug || Conf.App.Port != "3000" {
		t.Fatal("Config file is not initialized properly for production environment")
	}

	// testore test environment defaults
	InitConf("test")
}

func TestGoEnv(t *testing.T) {
	// GO_ENV should be used if "TITAN_ENV" is empty
	// if both are empty, should be env = dev
}
