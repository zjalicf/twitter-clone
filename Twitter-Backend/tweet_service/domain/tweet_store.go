package domain

import "context"

type TweetStore interface {
	GetFeedByUser(followings []string) ([]*Tweet, error)
	SaveImageRedis(imageBytes []byte) error
	GetAll(ctx context.Context) ([]Tweet, error)
	GetTweetsByUser(ctx context.Context, username string) ([]*Tweet, error)
	Post(ctx context.Context, tweet *Tweet) (*Tweet, error)
	Favorite(ctx context.Context, id string, username string) (int, error)
	GetLikesByTweet(ctx context.Context, tweetID string) ([]*Favorite, error)
}
