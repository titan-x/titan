package main

// Chat is a private or group chat.
type Chat struct {
	ID    uint64
	Users []User
}

// User is a mobile user.
type User struct {
	ID      uint32
	Devices []Device
	Name    string
	Picture []byte
}

// Device is an NBusy installed device.
type Device interface {
	Send()
}

// Android device.
type Android struct {
	GCMRegID string
	Phone    uint64
}

// iOS device.
type iOS struct {
	APNSDeviceToken string
	Phone           uint64
}

var chats = make(map[string]Chat)

// user -> id (user or chat id) -> message
// delivery status -> user
// read status -> user
