package main

import "crypto/tls"

// User is a mobile user.
type User struct {
	ID              uint32
	PhoneNumber     uint64
	GCMRegID        string
	APNSDeviceToken string
	Name            string
	Picture         []byte
	Conn            *tls.Conn
}

// Send sends given data to to a device using device specific infrastructure.
func (u *User) Send(data map[string]string) error {
	// note: not adding SendMessage/SendNotification/etc. like fine grained methods to keep this library more low level, those are in the mobile conn and queue types
	return nil
}
