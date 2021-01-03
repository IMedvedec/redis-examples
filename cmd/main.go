package main

import (
	"context"
	"log"

	"github.com/imedvedec/redis-examples/redis"
)

func main() {
	redisService, err := redis.NewService(":6379")
	if err != nil {
		log.Println(err)
		return
	}

	val, err := redisService.Ping(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(val)

	if err := redisService.Shutdown(context.Background()); err != nil {
		log.Println(err)
		return
	}
}
