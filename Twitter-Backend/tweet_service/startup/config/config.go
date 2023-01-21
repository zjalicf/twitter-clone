package config

import "os"

type Config struct {
	Port           string
	TweetDB        string
	JaegerAddress  string
	TweetCacheHost string
	TweetCachePort string
}

func NewConfig() *Config {
	return &Config{
		Port:           os.Getenv("TWEET_SERVICE_PORT"),
		TweetDB:        os.Getenv("TWEET_DB"),
		JaegerAddress:  os.Getenv("JAEGER_ADDRESS"),
		TweetCacheHost: os.Getenv("TWEET_CACHE_HOST"),
		TweetCachePort: os.Getenv("TWEET_CACHE_PORT"),
	}
}
