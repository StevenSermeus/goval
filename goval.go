package main

import (
	"os"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/logging"
	"github.com/StevenSermeus/goval/src/networking"
)

func main() {
	logging.Info.Println("Starting goval")
	logging.Info.Println("Loading config")
	config, err := config.LoadConfig()
	if err != nil {
		logging.Info.Println("Error loading config file, using default values")
		panic(err)
	}
	var cache = cache.Cache{}
	logging.Info.Println("Creating data directory")
	os.MkdirAll(config.DataDir, 0755)
	logging.Info.Println("Creating config directory")
	os.MkdirAll("./goval/config", 0755)
	go cache.CacheSizeManagement(uint64(config.MaxCacheSize))
	networking.Tcp(&cache, &config)
}
