package config

import "os"

type Config struct {
	Port          string
	TweetDB       string
	JaegerAddress string
}

func NewConfig() *Config {
	return &Config{
		Port:          os.Getenv("TWEET_SERVICE_PORT"),
		TweetDB:       os.Getenv("TWEET_DB"),
		JaegerAddress: os.Getenv("JAEGER_ADDRESS"),
	}
}
