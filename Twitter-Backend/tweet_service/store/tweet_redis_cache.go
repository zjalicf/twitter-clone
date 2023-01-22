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

func (a *TweetRedisCache) PostCacheData(key string, value *[]byte) error {
	log.Println("redis post")
	result := a.client.Set(key, *value, 10*time.Minute)
	log.Println(result.Err())
	log.Println(result.Result())
	if result.Err() != nil {
		log.Printf("redis set error: %s", result.Err())
		return result.Err()
	}

	return nil
}

func (a *TweetRedisCache) GetCachedValue(key string) (*[]byte, error) {
	result := a.client.Get(key)


	if result.Err() == nil {
		token, err := result.Bytes()
		if err != nil {
			return nil, err
		}
		return &token, nil
	}
	return nil, nil
}
