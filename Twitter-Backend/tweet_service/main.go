package main

import "Twitter-Backend/startup"
import cfg "Twitter-Backend/startup/config"

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
