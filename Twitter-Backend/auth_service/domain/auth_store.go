package domain

type AuthStore interface {
	Register(credentials *Credentials) error
	Login(credentials *Credentials) (string, error)
}
