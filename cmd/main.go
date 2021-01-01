package main

import (
	"github.com/imedvedec/redis-examples/redis"
)

func main() {
	sr := redis.NewService()

	sr.Ping()
	sr.Shutdown()
}
