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
	OwnerUsername string     `json:"owner_username"`
	Image         bool       `json:"image"`
	Advertisement bool       `json:"advertisement"`
}

type Favorite struct {
	TweetID  gocql.UUID `json:"tweet_id"`
	Username string     `json:"username"`
	ID       gocql.UUID `json:"id"`
}

type Retweet struct {
	TweetID  gocql.UUID `json:"tweet_id"`
	Username string     `json:"username"`
	ID       gocql.UUID `json:"id"`
}

type TweetID struct {
	ID string `json:"id"`
}

type Event struct {
	TweetID   string
	Type      string
	Timestamp int
}
