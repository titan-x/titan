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

// InMemDB is an in-memory database.
type InMemDB struct {
	InMemUserDB
}

// InMemUserDB is in-memory user database.
type InMemUserDB struct {
	ids    map[uint32]*User
	emails map[string]*User
}

// GetByID retrieves a user by ID.
func (db *InMemUserDB) GetByID(id uint32) (user *User, ok bool) {
	user, ok = db.ids[id]
	return
}

// GetByMail retrieves a user by e-mail address.
func (db *InMemUserDB) GetByMail(email string) (user *User, ok bool) {
	user, ok = db.emails[email]
	return
}
