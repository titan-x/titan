package main

import "os"

// GCM type describes the Google Cloud Messaging parameters as described here: https://developer.android.com/google/gcm/gs.html
type GCM struct {
	CCSEndpoint string
	SenderID    string
	APIKey      string
}

func config() {
	env := os.Getenv("GO_ENV")
	gcm := GCM{SenderID: os.Getenv("GCM_SENDER_ID"), APIKey: os.Getenv("GOOGLE_API_KEY")}
	if (env == "development") {
		gcm.CCSEndpoint = os.Getenv("GCM_CCS_ENDPOINT")
	} else {
		gcm.CCSEndpoint = os.Getenv("GCM_CCS_STAGING_ENDPOINT")
	}
}
