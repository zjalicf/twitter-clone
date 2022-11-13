package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"id" json:"id"`
	Firstname string             `bson:"firstName,omitempty" json:"firstName"`
	Lastname  string             `bson:"lastName,omitempty" json:"lastName"`
	Gender    Gender             `bson:"gender,omitempty" json:"gender"`
	Age       int                `bson:"age,omitempty" json:"age"`
	Residence string             `bson:"residence,omitempty" json:"residence"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	UserType  UserType           `bson:"userType" json:"userType"`

	CompanyName string `bson:"companyName,omitempty" json:"companyName"`
	Email       string `bson:"email,omitempty" json:"email"`
	Website     string `bson:"website,omitempty" json:"website"`
}

type Gender string

const (
	Male   = "Male"
	Female = "Female"
)

type UserType string

const (
	Regular  = "Regular"
	Business = "Business"
)
