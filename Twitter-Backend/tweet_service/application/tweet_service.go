package application

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"tweet_service/domain"
)

var (
	followServiceHost = os.Getenv("FOLLOW_SERVICE_HOST")
	followServicePort = os.Getenv("FOLLOW_SERVICE_PORT")
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

func (service *TweetService) GetTweetsByUser(username string) ([]*domain.Tweet, error) {
	return service.store.GetTweetsByUser(username)
}

func (service *TweetService) GetFeedByUser(token string) ([]*domain.Tweet, error) {

	followServiceEndpoint := fmt.Sprintf("http://%s:%s/followings", followServiceHost, followServicePort)
	followServiceRequest, _ := http.NewRequest("GET", followServiceEndpoint, nil)
	followServiceRequest.Header.Add("Authorization", token)
	responseFservice, _ := http.DefaultClient.Do(followServiceRequest)

	responseBodyBytes, err := io.ReadAll(responseFservice.Body)
	if err != nil {
		log.Printf("error in readAll: %s", err.Error())
		return nil, err
	}

	var followingsList []string
	err = json.Unmarshal(responseBodyBytes, &followingsList)
	if err != nil {
		log.Printf("error in unmarshal: %s", err.Error())
		return nil, err
	}

	userFeed, err := service.store.GetFeedByUser(followingsList)
	if err != nil {
		log.Printf("Error in getting feed by user: %s", err.Error())
		return nil, err
	}

	return userFeed, nil
}

func (service *TweetService) GetLikesByTweet(tweetID string) ([]*domain.Favorite, error) {
	return service.store.GetLikesByTweet(tweetID)
}

func (service *TweetService) Post(tweet *domain.Tweet, username string) (*domain.Tweet, error) {
	tweet.ID, _ = gocql.RandomUUID()
	tweet.CreatedAt = time.Now().Unix()
	tweet.Favorited = false
	tweet.FavoriteCount = 0
	tweet.Retweeted = false
	tweet.RetweetCount = 0
	tweet.Username = username

	return service.store.Post(tweet)
}

func (service *TweetService) Favorite(id *gocql.UUID) (int, error) {
	return service.store.Favorite(id)
}
