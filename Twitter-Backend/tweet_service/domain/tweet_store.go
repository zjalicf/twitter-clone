package domain

type TweetStore interface {
	//Get(id primitive.ObjectID) (*Tweet, error)
	GetAll() ([]Tweet, error)
	GetTweetsByUser(username string) ([]*Tweet, error)
	Post(tweet *Tweet) (*Tweet, error)
	//DeleteAll()
}
