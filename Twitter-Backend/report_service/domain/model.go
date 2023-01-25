package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID          primitive.ObjectID `bson:"_id"`
	TweetID     string             `json:"tweet_id" bson:"tweet_id"`
	Timestamp   int64              `json:"timestamp" bson:"timestamp"`
	LikeCount   int                `json:"like_count" bson:"like_count"`
	UnlikeCount int                `json:"unlike_count" bson:"unlike_count"`
	ViewCount   int                `json:"view_count" bson:"view_count"`
	TimeSpent   int                `json:"time_spended" bson:"time_spended"`
}

type Event struct {
	TweetID   string
	Type      string
	Timestamp int
}
