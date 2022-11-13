package domain


import "go.mongodb.org/mongo-driver/bson/primitive"
type AuthStore interface {
	Register(user *User, isBusiness bool) error
}
