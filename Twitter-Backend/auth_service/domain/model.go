package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"id" json:"id"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Gender    Gender             `bson:"gender" json:"gender"`
	Age       int                `bson:"age" json:"age"`
	Residence string             `bson:"residence" json:"residence"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	UserType  UserType           `bson:"userType" json:"userType"`

	CompanyName string `bson:"companyName,omitempty" json:"companyName,omitempty"`
	Email       string `bson:"email,omitempty" json:"email,omitempty"`
	WebSite     string `bson:"webSite,omitempty" json:"webSite,omitempty"`
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

type Credentials struct {
	Username	string	`bson:"username" json:"username"`
	Password	string	`bson:"password" json:"password"`
}