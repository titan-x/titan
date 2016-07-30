package titan

// DB wraps all database related functions.
type DB interface {
	UserDB
}

// UserDB presists user information in database.
type UserDB interface {
	Seed() error
	GetByID(id string) (*User, bool)
	GetByMail(mail string) (*User, bool)
	SaveUser(u *User) error
}
