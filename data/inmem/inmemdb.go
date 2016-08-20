package inmem

import (
	"strconv"

	"github.com/titan-x/titan/models"
)

// DB is an in-memory database.
type DB struct {
	UserDB
}

// UserDB is in-memory user database.
type UserDB struct {
	ids    map[string]*models.User
	emails map[string]*models.User
}

// NewDB creates a new in-memory database.
func NewDB() DB {
	return DB{
		UserDB: UserDB{
			ids:    make(map[string]*models.User),
			emails: make(map[string]*models.User),
		},
	}
}

// Seed seeds database with essential data.
func (db UserDB) Seed(overwrite bool) error {
	return nil
}

// GetByID retrieves a user by ID.
func (db UserDB) GetByID(id string) (u *models.User, ok bool) {
	u, ok = db.ids[id]
	return
}

// GetByMail retrieves a user by e-mail address.
func (db UserDB) GetByMail(email string) (u *models.User, ok bool) {
	u, ok = db.emails[email]
	return
}

// SaveUser save or updates a user object in the database.
func (db UserDB) SaveUser(u *models.User) error {
	if u.ID == "" {
		u.ID = strconv.Itoa(len(db.ids) + 1)
	}

	db.ids[u.ID] = u
	db.emails[u.Email] = u
	return nil
}
