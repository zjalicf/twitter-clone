package domain

import (
	"github.com/gocql/gocql"
)

type Tweet struct {
	ID            gocql.UUID `json:"id"`
	Text          string     `json:"text"`
	CreatedAt     int64      `json:"created_on"`
	Favorited     bool       `json:"favorited"`
	FavoriteCount int        `json:"favorite_count"`
	Retweeted     bool       `json:"retweeted"`
	RetweetCount  int        `json:"retweet_count"`
	Username      string     `json:"username"`
	Image         bool       `json:"bool"`
}

type Favorite struct {
	TweetID  gocql.UUID `json:"tweet_id"`
	Username string     `json:"username"`
	ID       gocql.UUID `json:"id"`
}

type TweetID struct {
	ID string `json:"id"`
}
