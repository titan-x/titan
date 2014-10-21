package main

import (
	"fmt"
	"os"
)

const (
	gcmSenderID           = "218602439235"
	gcmCcsEndpoint        = "gcm.googleapis.com:5235"
	gcmPrepodSenderID     = ""
	gcmCcsPreprodEndpoint = "gcm-preprod.googleapis.com:5236"

	// GO execution envrinment variables
	goEnv   = "GO_ENV"
	goDebug = "GO_DEBUG"

	// possible GO_ENV values
	dev     = "development"
	test    = "test"
	staging = "staging"
	prod    = "production"

	// Google environment variables
	googAPIKey        = "GOOGLE_API_KEY"
	googPreprodAPIKey = "GOOGLE_PREPROD_API_KEY"

	// GCM environment variables
	gcmRegID = "GCM_REG_ID"
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
	// One of the following: development, test, staging, production
	Env string
	// Enables verbose logging to stdout
	Debug bool
}

// GCM describes the Google Cloud Messaging parameters as described here: https://developer.android.com/google/gcm/gs.html
type GCM struct {
	CCSEndpoint string
	SenderID    string
	APIKey      string
	RegID       string
}

// GetConfig returns a singleton instance of the application configuration.
func GetConfig() Config {
	if initialized {
		return config
	}

	env := os.Getenv(goEnv)
	if env == "" {
		env = dev
	}

	debug := os.Getenv(goDebug) != ""

	app := App{Env: env, Debug: debug}

	gcm := GCM{CCSEndpoint: gcmCcsEndpoint, SenderID: gcmSenderID, APIKey: os.Getenv(googAPIKey), RegID: os.Getenv(gcmRegID)}
	if env != prod && os.Getenv(googPreprodAPIKey) != "" {
		// todo: use preprod specific endpoint, sender ID, and API key from a separate app (i.e. nbusy-preprod)
	}

	config = Config{App: app, GCM: gcm}
	if debug {
		fmt.Printf("Config: %+v\n", config)
	}

	initialized = true
	return config
}
