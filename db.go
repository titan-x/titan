package devastator

// DB wraps all database related functions.
type DB interface {
	UserDB
}

// UserDB presists user information in database.
type UserDB interface {
	GetByID(id uint32) (*User, bool)
	GetByMail(mail string) (*User, bool)
}

// InMemoryDB is an in-memory database.
type InMemoryDB struct {
	InMemoryUserDB
}

// InMemoryUserDB is in-memory user database.
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
