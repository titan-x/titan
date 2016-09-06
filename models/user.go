package models

import "time"

// User profile.
type User struct {
	ID              string
	Registered      time.Time
	Email           string
	PhoneNumber     string
	GCMRegID        string
	APNSDeviceToken string
	Name            string
	Picture         []byte
	JWTToken        string
}
