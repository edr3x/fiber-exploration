package config

import (
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func RedisConnect() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if Redis == nil {
		panic("Redis connect error")
	}

	log.Println("Connected to Redis")
}
