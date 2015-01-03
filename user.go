package main

// User is a mobile user.
type User struct {
	ID          uint32
	PhoneNumber uint64
	GCMRegID    string
	Name        string
	Picture     []byte
}

// Send sends given data to to a device using device specific infrastructure.
func (u *User) Send(data map[string]string) error {
	// note: not adding SendMessage/SendNotification/etc. like fine grained methods to keep this library more low level
	return nil
}

// user -> id (user or chat id) -> message
// delivery status -> user
// read status -> user
