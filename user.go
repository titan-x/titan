package devastator

import "github.com/nbusy/neptulon"

// User encapsulates connected user information.
type User struct {
	ID              uint32
	Email           string
	PhoneNumber     uint64
	GCMRegID        string
	APNSDeviceToken string
	Name            string
	Picture         []byte
	Conn            *neptulon.Conn
	Cert            []byte

	// MsgQueue may contain request, response, or notification messages.
	MsgQueue []interface{}
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
