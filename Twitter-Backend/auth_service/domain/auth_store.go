package domain

type AuthStore interface {
	Register(user *User, isBusiness bool) error
}