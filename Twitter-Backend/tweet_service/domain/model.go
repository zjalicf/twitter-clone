package domain

import (
	"github.com/gocql/gocql"
)

type Tweet struct {
	ID   gocql.UUID `bson:"id" json:"id"`
	Text string     `bson:"text" json:"text"`
	//Image     string             `bson:"image,omitempty" json:"image,omitempty"`
	//CreatedOn time.Time          `bson:"created_on" json:"created_on"`
}
