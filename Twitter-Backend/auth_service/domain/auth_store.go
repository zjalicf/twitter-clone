package domain

type AuthStore interface {
	GetAll() ([]*User, error)
	Register(user *Credentials) error
	GetOneUser(username string) (*User, error)
}
