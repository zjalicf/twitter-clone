package store

import (
	"github.com/go-redis/redis"
	"log"
	"time"
	"tweet_service/domain"
)

type TweetRedisCache struct {
	client *redis.Client
}

func NewTweetRedisCache(client *redis.Client) domain.TweetCache {
	return &TweetRedisCache{
		client: client,
	}
}

func (a *TweetRedisCache) PostCacheData(key string, value string) error {
	log.Println("redis post")
	result := a.client.Set(key, value, 10*time.Minute)
	log.Println(result.Err())
	log.Println(result.Result())
	if result.Err() != nil {
		log.Printf("redis set error: %s", result.Err())
		return result.Err()
	}

	return nil
}

func (a *TweetRedisCache) GetCachedValue(key string) (string, error) {
	result := a.client.Get(key)
	token, err := result.Result()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return token, nil
}
