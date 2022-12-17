package domain

import (
	"github.com/gocql/gocql"
)

type Tweet struct {
	ID   gocql.UUID `json:"id"`
	Text string     `json:"text"`
	//Image     string         `bson:"image,omitempty" json:"image,omitempty"`
	CreatedAt     int64  `json:"created_on"`
	Favorited     bool   `json:"favorited"`
	FavoriteCount int    `json:"favorite_count"`
	Retweeted     bool   `json:"retweeted"`
	RetweetCount  int    `json:"retweet_count"`
	Username      string `json:"username"`
}

type Favorite struct {
	TweetID  gocql.UUID `json:"id"`
	Username string     `json:"username"`
}

type TweetID struct {
	ID string `json:"id"`
}
