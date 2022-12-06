package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthStore interface {
	GetAll() ([]*User, error)
	Register(user *Credentials) error
	GetOneUser(username string) (*User, error)
	GetOneUserByID(id primitive.ObjectID) *User
	ChangePassword(user *User) error
}
