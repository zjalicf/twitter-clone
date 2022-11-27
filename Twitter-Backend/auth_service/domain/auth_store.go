package domain

type AuthStore interface {
	Login(credentials *Credentials) (string, error)
}
