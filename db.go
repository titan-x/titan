package titan

// DB wraps all database related functions.
type DB interface {
	UserDB
}

// UserDB presists user information in database.
type UserDB interface {
	Seed(overwrite bool) error
	GetByID(id string) (u *User, ok bool)
	GetByMail(email string) (u *User, ok bool)
	SaveUser(u *User) error
}
