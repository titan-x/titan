package devastator

// User encapsulates user information.
type User struct {
	ID              string
	Email           string
	PhoneNumber     uint64
	GCMRegID        string
	APNSDeviceToken string
	Name            string
	Picture         []byte
	Cert            []byte
}

// Send sends given data to to a device using device specific infrastructure.
// todo: not adding SendMessage/SendNotification/etc. like fine grained methods to keep this library more low level, those are in the mobile conn and queue types
func (u *User) Send(msg interface{}) error {
	return nil
}

// Queue queues a message to be sent to a user as soon as possible.
func (u *User) Queue(msg interface{}) error {
	return nil
}
