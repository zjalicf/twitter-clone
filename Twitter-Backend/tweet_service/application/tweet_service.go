package application

import (
	"context"
	"github.com/gocql/gocql"
	"go.opentelemetry.io/otel/trace"
	"time"
	"tweet_service/domain"
)

type TweetService struct {
	store  domain.TweetStore
	tracer trace.Tracer
}

func NewTweetService(store domain.TweetStore, tracer trace.Tracer) *TweetService {
	return &TweetService{
		store:  store,
		tracer: tracer,
	}
}

func (service *TweetService) GetAll(ctx context.Context) ([]domain.Tweet, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.GetAll")
	defer span.End()

	return service.store.GetAll()
}

func (service *TweetService) GetTweetsByUser(ctx context.Context, username string) ([]*domain.Tweet, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.GetTweetsByUser")
	defer span.End()

	return service.store.GetTweetsByUser(username)
}

func (service *TweetService) GetLikesByTweet(ctx context.Context, tweetID string) ([]*domain.Favorite, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.GetLikesByTweet")
	defer span.End()

	return service.store.GetLikesByTweet(tweetID)
}

func (service *TweetService) Post(ctx context.Context, tweet *domain.Tweet, username string) (*domain.Tweet, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.Post")
	defer span.End()

	tweet.ID, _ = gocql.RandomUUID()
	tweet.CreatedAt = time.Now().Unix()
	tweet.Favorited = false
	tweet.FavoriteCount = 0
	tweet.Retweeted = false
	tweet.RetweetCount = 0
	tweet.Username = username

	return service.store.Post(tweet)
}

func (service *TweetService) Favorite(ctx context.Context, id string, username string) (int, error) {
	ctx, span := service.tracer.Start(ctx, "TweetService.Favorite")
	defer span.End()

	return service.store.Favorite(id, username)
}
