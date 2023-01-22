package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	TweetID     primitive.ObjectID `bson:"tweet_id"`
	LikeCount   int                `bson:"like_count"`
	UnlikeCount int                `bson:"unlike_count"`
	ViewCount   int                `bson:"view_count"`
}

type Event struct {
	TweetID   primitive.ObjectID
	Type      string
	Timestamp int
}
