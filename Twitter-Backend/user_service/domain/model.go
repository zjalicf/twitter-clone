package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"id" json:"id"`
	Firstname string             `bson:"firstname,omitempty" json:"firstname"`
	Lastname  string             `bson:"lastname,omitempty" json:"lastname"`
	Gender    Gender             `bson:"gender,omitempty" json:"gender"`
	Age       int                `bson:"age,omitempty" json:"age"`
	Address   string             `bson:"address,omitempty" json:"lastname"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`

	CompanyName string `bson:"companyName,omitempty" json:"companyName"`
	Email       string `bson:"email,omitempty" json:"email"`
	Website     string `bson:"website,omitempty" json:"website"`
}

type Gender string

const (
	Male   = "Male"
	Female = "Female"
)
