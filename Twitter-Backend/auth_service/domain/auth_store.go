package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthStore interface {
	Register(user newUser) (*User, error)
	Login(credentials Credentials) (string, error)
}
