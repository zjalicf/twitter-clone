package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type TweetStore interface {
	Get(id primitive.ObjectID) (*Tweet, error)
	GetAll() ([]*Tweet, error)
	Post(tweet *Tweet) error
	DeleteAll()
}
