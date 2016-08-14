package titan

// DB wraps all database related functions.
type DB interface {
	UserDB
}

// UserDB presists user information in database.
type UserDB interface {
	Seed(overwrite bool) error
	GetByID(id string) (*User, bool)
	GetByMail(mail string) (*User, bool)
	SaveUser(u *User) error
}
