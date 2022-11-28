package domain

type AuthStore interface {
	GetOneUser(username string) (*User, error)
}
