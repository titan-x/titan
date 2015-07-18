package devastator

// UserDB presists user information in database.
type UserDB interface {
	GetByID(id uint32) (*User, bool)
	GetByMail(mail string) (*User, bool)
}

// InMemoryUserDB is an in-memory user DB which is useful for caching user objects retrieved from database, or testing.
type InMemoryUserDB struct {
	ids   map[uint32]*User
	mails map[string]*User
}

// GetByID retrieves a user by ID.
func (u *InMemoryUserDB) GetByID(id uint32) (*User, bool) {
	return nil, false
}

// GetByMail retrieves a user by e-mail address.
func (u *InMemoryUserDB) GetByMail(mail string) (*User, bool) {
	return nil, false
}
