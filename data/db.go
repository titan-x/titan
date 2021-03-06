package data

import "github.com/titan-x/titan/models"

// DB wraps all database related functions.
type DB interface {
	UserDB
}

// UserDB presists user information in database.
type UserDB interface {
	Seed(overwrite bool, jwtPass string) error
	GetByID(id string) (u *models.User, ok bool)
	GetByEmail(email string) (u *models.User, ok bool)
	SaveUser(u *models.User) error
}
