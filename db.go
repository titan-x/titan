package titan

// DB wraps all database related functions.
type DB interface {
	UserDB
}

// todo: db.GetByMail is not good, either use udb.GetByMail -or- db.Users.GetByMail

// UserDB presists user information in database.
type UserDB interface {
	GetByID(id string) (*User, bool)
	GetByMail(mail string) (*User, bool)
	SaveUser(u *User) error
}
