package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tweet struct {
	ID        primitive.ObjectID `bson:"id"`
	Text      string             `bson:"text"`
	Image     string             `bson:"image"`
	CreatedOn time.Time          `bson:"created_on"`
}
