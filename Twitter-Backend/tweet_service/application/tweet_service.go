package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/sony/gobreaker"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	"go.opentelemetry.io/otel/trace"
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

//eesa

type TweetService struct {
	store        domain.TweetStore
	tracer       trace.Tracer
	cache        domain.TweetCache
	cb           *gobreaker.CircuitBreaker
	orchestrator *CreateEventOrchestrator
}

func NewTweetService(store domain.TweetStore, cache domain.TweetCache, tracer trace.Tracer, orchestrator *CreateEventOrchestrator) *TweetService {
	return &TweetService{
		store:        store,
		cache:        cache,
		cb:           CircuitBreaker(),
		orchestrator: orchestrator,
		tracer:       tracer,
	}
}

func (service *TweetService) GetAll(ctx context.Context) ([]domain.Tweet, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.GetAll")
	defer span.End()

	return service.store.GetAll(ctx)
}

func (service *TweetService) GetTweetsByUser(ctx context.Context, username string) ([]*domain.Tweet, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.GetTweetsByUser")
	defer span.End()

	return service.store.GetTweetsByUser(ctx, username)
}

func (service *TweetService) GetFeedByUser(token string) ([]*domain.Tweet, error) {

	followServiceEndpoint := fmt.Sprintf("http://%s:%s/followings", followServiceHost, followServicePort)
	followServiceRequest, _ := http.NewRequest("GET", followServiceEndpoint, nil)
	followServiceRequest.Header.Add("Authorization", token)
	bodyBytes, err := service.cb.Execute(func() (interface{}, error) {

		responseFservice, err := http.DefaultClient.Do(followServiceRequest)
		if err != nil {
			return nil, fmt.Errorf("FollowServiceError")
		}

		defer responseFservice.Body.Close()

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

		return followingsList, nil
	})

	if err != nil {
		return nil, err
	}

	userFeed, err := service.store.GetFeedByUser(bodyBytes.([]string))
	if err != nil {
		log.Printf("Error in getting feed by user: %s", err.Error())
		return nil, err
	}

	return userFeed, nil
}

func (service *TweetService) saveImage(tweetID gocql.UUID, imageBytes []byte) error {
	return service.store.SaveImage(tweetID, imageBytes)
}

func (service *TweetService) GetLikesByTweet(ctx context.Context, tweetID string) ([]*domain.Favorite, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.GetLikesByTweet")
	defer span.End()

	return service.store.GetLikesByTweet(ctx, tweetID)
}

func (service *TweetService) Post(ctx context.Context, tweet *domain.Tweet, username string, image *[]byte) (*domain.Tweet, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.Post")
	defer span.End()

	log.Printf("Ad je : %t", tweet.Advertisement)

	tweet.ID, _ = gocql.RandomUUID()

	tweet.Image = false
	if len(*image) != 0 {
		log.Printf("USLO U SLIKU")
		err := service.saveImage(tweet.ID, *image)
		if err != nil {
			return nil, err
		}

		err = service.cache.PostCacheData(tweet.ID.String(), image)
		if err != nil {
			return nil, err
		}
		tweet.Image = true
	}
	tweet.CreatedAt = time.Now().Unix()
	tweet.Favorited = false
	tweet.FavoriteCount = 0
	tweet.Retweeted = false
	tweet.RetweetCount = 0
	tweet.Username = username

	return service.store.Post(ctx, tweet)
}

func (service *TweetService) Favorite(ctx context.Context, id string, username string, isAd bool) (int, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.Favorite")
	defer span.End()

	status, err := service.store.Favorite(ctx, id, username)
	if err != nil {
		log.Printf("Error in tweet_service Favorite(): %s", err.Error())
		return status, err
	}

	if isAd {
		event := events.Event{
			TweetID:   id,
			Type:      "",
			Timestamp: time.Now().Unix(),
			Timespent: 0,
		}
		if status == 200 {
			event.Type = "Unliked"
		} else {
			event.Type = "Liked"
		}
		err = service.orchestrator.Start(ctx, event)
		if err != nil {
			return status, err
		}
	}

	return status, nil
}

func (service *TweetService) TimeSpentOnAd(ctx context.Context, timespent *domain.Timespent) error {
	ctx, span := service.tracer.Start(ctx, "TweetService.TimeSpentOnAd")
	defer span.End()
	log.Printf("TIMESPENT IN tweet_service.TimeSpentOnAd: %s", timespent)
	event := events.Event{
		TweetID:   timespent.TweetID,
		Type:      "Timespent",
		Timestamp: time.Now().Unix(),
		Timespent: timespent.Timespent,
	}

	err := service.orchestrator.Start(ctx, event)
	if err != nil {
		log.Printf("Error in starting orchestrator in TweetService.TimeSpentOnAd: %s", err.Error())
		return err
	}

	return nil
}

func (service *TweetService) GetTweetImage(id string) (*[]byte, error) {
	cachedImage, _ := service.cache.GetCachedValue(id)
	//if err != nil {
	//	log.Printf("GET REDIS ERR: %s", err.Error())
	//	return nil, err
	//}

	if cachedImage != nil {

		return cachedImage, nil
	}

	image, err := service.store.GetTweetImage(context.TODO(), id)
	if err != nil {
		log.Printf("CASSANDRA ERR: %s", err.Error())
		return nil, err
	}

	err = service.cache.PostCacheData(id, &image)
	if err != nil {
		log.Printf("POST REDIS ERR: %s", err.Error())
		return nil, err
	}
	return &image, nil
}

func (service *TweetService) Retweet(id string, username string) (int, error) {
	return service.store.Retweet(id, username)
}

func CircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "cb",
			MaxRequests: 1,
			Timeout:     time.Millisecond,
			Interval:    0,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 3
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker '%s' changed from '%s' to '%s'\n", name, from, to)
			},
		},
	)
}
