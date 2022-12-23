package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthStore interface {
	GetAll() ([]*Credentials, error)
	Register(user *Credentials) error
	GetOneUser(username string) (*Credentials, error)
	GetOneUserByID(id primitive.ObjectID) *Credentials
	UpdateUser(user *Credentials) error
}
