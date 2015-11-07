package titan

// User profile.
type User struct {
	ID              string
	Email           string
	PhoneNumber     uint64
	GCMRegID        string
	APNSDeviceToken string
	Name            string
	Picture         []byte
	Cert            []byte // PEM encoded X.509 client-certificate
	Key             []byte // PEM encoded X.509 private key
}
