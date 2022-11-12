package main

import (
	"tweet-service/startup"
	cfg "tweet-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
