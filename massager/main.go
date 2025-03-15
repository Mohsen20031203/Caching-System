package main

import (
	"chach/massager/api"
	"chach/massager/config"
	"chach/massager/db"

	"github.com/go-redis/redis/v8"
)

func redisChack(addres string) *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr: addres,
	})
	return rds
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		return
	}

	rdb := redisChack("localhost:6379")

	db, err := db.NewStorege(config)
	if err != nil {
		return
	}

	server, err := api.NewServer(db, &config, rdb)
	if err != nil {
		return
	}
	server.Router.Run(":5436")
}
