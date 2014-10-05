package main

import (
	"os"
	"fmt"
)

const (
	gcmSenderID           = "218602439235"
	gcmCcsEndpoint        = "gcm.googleapis.com:5235"
	gcmCcsStagingEndpoint = "gcm-preprod.googleapis.com:5236"
)

var config Config
var initialized bool

// Config describes the global configuration for the NBusy server.
type Config struct {
	App App
	GCM GCM
}

// App contains the global application variables.
type App struct {
	Env   string // development, test, staging, production
	Debug bool
}

// GCM describes the Google Cloud Messaging parameters as described here: https://developer.android.com/google/gcm/gs.html
type GCM struct {
	CCSEndpoint string
	SenderID    string
	APIKey      string
}

// Config returns a singleton instance of the application configuration.
func GetConfig() Config {
	if initialized {
		return config
	}

	env := os.Getenv("GO_ENV")
	if (env == "") {
		env = "development"
	}
	debug := os.Getenv("GO_DEBUG") != ""
	app := App{Env: env, Debug: debug}

	gcm := GCM{CCSEndpoint: gcmCcsEndpoint, SenderID: gcmSenderID, APIKey: os.Getenv("GOOGLE_API_KEY")}
	if env != "production" {
		// todo: use staging specific endpoint, sender ID, and API key (i.e. nbusy-test)
	}

	config = Config{App: app, GCM: gcm}
	initialized = true
	if (debug) {
		fmt.Printf("Config: %+v\n", config)
	}
	return config
}
