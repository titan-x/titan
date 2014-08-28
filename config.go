package main

import "os"

const (
	gcmCcsEndpoint        = "gcm.googleapis.com:5235"
	gcmCcsStagingEndpoint = "gcm-staging.googleapis.com:5236"
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
	Env string
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

	app := App{Env: os.Getenv("GO_ENV")}

	gcm := GCM{SenderID: os.Getenv("GCM_SENDER_ID"), APIKey: os.Getenv("GOOGLE_API_KEY")}
	if app.Env == "development" {
		gcm.CCSEndpoint = gcmCcsStagingEndpoint
	} else {
		gcm.CCSEndpoint = gcmCcsEndpoint
	}

	config = Config{App: app, GCM: gcm}
	initialized = true
	return config
}
