package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Firstname string             `bson:"firstName,omitempty" json:"firstName,omitempty"`
	Lastname  string             `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Gender    Gender             `bson:"gender,omitempty" json:"gender,omitempty"`
	Age       int                `bson:"age,omitempty" json:"age,omitempty"`
	Residence string             `bson:"residence,omitempty" json:"residence,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	UserType  UserType           `bson:"userType" json:"userType"`

	CompanyName string `bson:"companyName,omitempty" json:"companyName,omitempty"`
	Website     string `bson:"website,omitempty" json:"website,omitempty"`
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
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	UserType UserType           `bson:"userType" json:"userType"`
}

type Claims struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Username string             `json:"username"`
	Role     UserType           `json:"userType"`
	jwt.RegisteredClaims
}

type RegisterRecoverVerification struct {
	UserToken string `json:"user_token"`
	MailToken string `json:"mail_token"`
}

type ResendVerificationRequest struct {
	UserToken string `json:"user_token"`
	UserMail  string `json:"user_mail"`
}

type ResetPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	RepeatedNew string `json:"repeated_new"`
}

type RecoverPasswordRequest struct {
	UserID      string `json:"id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	RepeatedNew string `json:"repeated_new"`
}
