package domain

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"regexp"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Firstname string             `bson:"firstName,omitempty" json:"firstName,omitempty" validate:"regular"`
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
	UserID    primitive.ObjectID `json:"user_id"`
	Role      UserType           `json:"userType"`
	ExpiresAt time.Time          `json:"expires_at"`
}

type RegisterValidation struct {
	UserToken string `json:"user_token"`
	MailToken string `json:"mail_token"`
}

func (user *User) ValidateRegular() error {
	validate := validator.New()

	err := validate.RegisterValidation("regular", regularField)
	if err != nil {
		return err
	}

	return validate.Struct(user)
}

func regularField(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("[-_a-zA-Z]*")
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

func (user *User) FromJSON(reader io.Reader) error {
	d := json.NewDecoder(reader)
	return d.Decode(user)
}
