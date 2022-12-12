package application

import (
	"github.com/gocql/gocql"
	"time"
	"tweet_service/domain"
)

type TweetService struct {
	store domain.TweetStore
}

func NewTweetService(store domain.TweetStore) *TweetService {
	return &TweetService{
		store: store,
	}
}

//func (service *TweetService) Get(id primitive.ObjectID) (*domain.Tweet, error) {
//	return service.store.Get(id)
//}

func (service *TweetService) GetAll() ([]domain.Tweet, error) {
	return service.store.GetAll()
}

func (service *TweetService) GetTweetsByUser(userID string) ([]*domain.Tweet, error) {
	return service.store.GetTweetsByUser(userID)
}

func (service *TweetService) Post(tweet *domain.Tweet, userID string) (*domain.Tweet, error) {
	tweet.ID, _ = gocql.RandomUUID()
	tweet.CreatedAt = time.Now().Unix()
	tweet.Favorited = false
	tweet.FavoriteCount = 0
	tweet.Retweeted = false
	tweet.RetweetCount = 0
	tweet.UserID = userID

	return service.store.Post(tweet)
}
