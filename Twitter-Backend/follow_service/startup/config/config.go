package config

import "os"

type Config struct {
	Port         string
	FollowDBHost string
	FollowDBPort string
	FollowDBUser string
	FollowDBPass string
}

func NewConfig() *Config {
	return &Config{
		Port:         os.Getenv("FOLLOW_SERVICE_PORT"),
		FollowDBHost: os.Getenv("FOLLOW_DB_HOST"),
		FollowDBPort: os.Getenv("FOLLOW_DB_PORT"),
		FollowDBUser: os.Getenv("FOLLOW_DB_USER"),
		FollowDBPass: os.Getenv("FOLLOW_DB_PASS"),
	}
}
