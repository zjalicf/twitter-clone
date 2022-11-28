package domain

type AuthStore interface {
	Register(user *Credentials) error
	GetOneUser(username string) (*User, error)
}
