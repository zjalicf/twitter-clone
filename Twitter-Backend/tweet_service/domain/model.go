package domain

import (
	"github.com/gocql/gocql"
)

type Tweet struct {
	ID   gocql.UUID `json:"id"`
	Text string     `json:"text"`
	//Image     string         `bson:"image,omitempty" json:"image,omitempty"`
	CreatedAt     int64      `json:"created_on"`
	Favorited     bool       `json:"favorited"`
	FavoriteCount int        `json:"favorite_count"`
	Retweeted     bool       `json:"retweeted"`
	RetweetCount  int        `json:"retweet_count"`
	UserID        gocql.UUID `json:"user_id"`
}
