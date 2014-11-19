package main

import (
	"log"
	"os"
)

const (
	// NBusy server envrinment variables
	nbusyEnv   = "NBUSY_ENV"
	nbusyDebug = "NBUSY_DEBUG"

	// possible NBUSY_ENV values
	dev     = "development"
	test    = "test"
	staging = "staging"
	prod    = "production"

	// GCM environment variables
	gcmSenderID = "GCM_SENDER_ID"
	gcmCcsHost  = "GCM_CCS_HOST"

	// Google environment variables
	googleAPIKey = "GOOGLE_API_KEY"
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
	CCSHost  string
	SenderID string
	APIKey   string
}

// GetConfig returns a singleton instance of the application configuration.
func GetConfig() Config {
	if initialized {
		return config
	}

	debug := os.Getenv(nbusyDebug) != ""
	env := os.Getenv(nbusyEnv)
	if env == "" {
		env = dev
	}

	app := App{Env: env, Debug: debug}
	gcm := GCM{CCSHost: gcmCcsHost, SenderID: gcmSenderID, APIKey: os.Getenv(googleAPIKey)}
	config = Config{App: app, GCM: gcm}

	if debug {
		log.Printf("Config: %+v\n", config)
	}

	initialized = true
	return config
}
