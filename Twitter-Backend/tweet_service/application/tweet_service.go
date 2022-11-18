package application

import (
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

func (service *TweetService) GetAll() ([]*domain.Tweet, error) {
	return service.store.GetAll()
}

//func (service *TweetService) Post(tweet *domain.Tweet) error {
//	tweet.ID = primitive.NewObjectID()
//	tweet.CreatedOn = time.Now()
//	return service.store.Post(tweet)
//}
