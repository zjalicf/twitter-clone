package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	UserType UserType           `bson:"userType" json:"userType"`
}

type UserType string

const (
	Regular  = "Regular"
	Business = "Business"
)

type Claims struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Username string             `json:"username"`
	Role     UserType           `json:"userType"`
	jwt.RegisteredClaims
}
