package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tweet struct {
	ID        primitive.ObjectID `bson:"id" json:"id"`
	Text      string             `bson:"text" json:"text"`
	Image     string             `bson:"image,omitempty" json:"image,omitempty"`
	CreatedOn time.Time          `bson:"created_on" json:"created_on"`
}
