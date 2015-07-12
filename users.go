package devastator

// Users is presists user information in database.
type Users struct {
	ids   map[uint32]*User
	mails map[string]*User
}

// GetByID retrieves a user by ID.
func (u *Users) GetByID(id uint32) (*User, bool) {
	return nil, false
}

// GetByMail retrieves a user by e-mail address.
func (u *Users) GetByMail(mail string) (*User, bool) {
	return nil, false
}
