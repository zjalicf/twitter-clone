package config

import "os"

type Config struct {
	Port       string
	AuthDBHost string
	AuthDBPORT string
}

func NewConfig() *Config {
	return &Config{
		Port:       os.Getenv("AUTH_SERVICE_PORT"),
		AuthDBHost: os.Getenv("AUTH_DB_HOST"),
		AuthDBPORT: os.Getenv("AUTH_DB_PORT"),
	}
}
