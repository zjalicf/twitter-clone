package domain

type TweetStore interface {
	//Get(id primitive.ObjectID) (*Tweet, error)
	GetAll() ([]Tweet, error)
	GetTweetsByUser(username string) ([]*Tweet, error)
	GetFeedByUser(followings []string) ([]*Tweet, error)
	Post(tweet *Tweet) (*Tweet, error)
	Favorite(id string, username string) (int, error)
	GetLikesByTweet(tweetID string) ([]*Favorite, error)
	SaveImageRedis(imageBytes []byte) error
}
