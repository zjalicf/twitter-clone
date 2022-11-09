package domain

type TweetRepo interface {
	GetAll() Tweets
	PostTweet(tweet *Tweet) Tweet
	PutTweet(tweet *Tweet, id int) error
	DeleteTweet(id int) error
}
