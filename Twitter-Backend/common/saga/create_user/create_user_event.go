package create_user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID
	Firstname  string
	Lastname   string
	Gender     Gender
	Age        int
	Residence  string
	Email      string
	Username   string
	Password   string
	UserType   UserType
	Visibility bool

	CompanyName string `bson:"companyName,omitempty" json:"companyName,omitempty" validation:"onlyCharAndNum"`
	Website     string `bson:"website,omitempty" json:"website,omitempty" validate:"onlyCharAndNum"`
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

type CreateUserCommandType int8

const (
	UpdateAuth CreateUserCommandType = iota
	UpdateUsers
	UpdateGraph
	UnknownCommand
)

type CreateUserCommand struct {
	User User
	Type CreateUserCommandType
}

type CreateUserReplyType int8

const (
	AuthUpdated CreateUserReplyType = iota
	UsersUpdated
	GraphUpdated
	UnknownReply
)

type CreateUserReply struct {
	User User
	Type CreateUserReplyType
}

//aaaa
//aaaa
