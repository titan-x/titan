package titan

import "strconv"

// InMemDB is an in-memory database.
type InMemDB struct {
	InMemUserDB
}

// InMemUserDB is in-memory user database.
type InMemUserDB struct {
	ids    map[string]*User
	emails map[string]*User
}

// NewInMemDB creates a new in-memory database.
func NewInMemDB() InMemDB {
	return InMemDB{
		InMemUserDB: InMemUserDB{
			ids:    make(map[string]*User),
			emails: make(map[string]*User),
		},
	}
}

// Seed seeds database with essential data.
func (db InMemUserDB) Seed(overwrite bool) error {
	return nil
}

// GetByID retrieves a user by ID.
func (db InMemUserDB) GetByID(id string) (u *User, ok bool) {
	u, ok = db.ids[id]
	return
}

// GetByMail retrieves a user by e-mail address.
func (db InMemUserDB) GetByMail(email string) (u *User, ok bool) {
	u, ok = db.emails[email]
	return
}

// SaveUser save or updates a user object in the database.
func (db InMemUserDB) SaveUser(u *User) error {
	if u.ID == "" {
		u.ID = strconv.Itoa(len(db.ids) + 1)
	}

	db.ids[u.ID] = u
	db.emails[u.Email] = u
	return nil
}
