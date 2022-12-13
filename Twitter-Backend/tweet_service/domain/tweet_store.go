package domain

import "github.com/gocql/gocql"

type TweetStore interface {
	//Get(id primitive.ObjectID) (*Tweet, error)
	GetAll() ([]Tweet, error)
	GetTweetsByUser(username string) ([]*Tweet, error)
	Post(tweet *Tweet) (*Tweet, error)
	Favorite(id *gocql.UUID) (int, error)
	//DeleteAll()
}
