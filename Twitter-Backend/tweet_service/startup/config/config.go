package config

import "os"

type Config struct {
	Port        string
	TweetDBHost string
	TweetDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port:        os.Getenv("TWEET_SERVICE_PORT"),
		TweetDBHost: os.Getenv("TWEET_DB_HOST"),
		TweetDBPort: os.Getenv("TWEET_DB_PORT"),
	}
}
