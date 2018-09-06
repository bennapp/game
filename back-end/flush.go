package main

import "github.com/go-redis/redis"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	redisClient.FlushAll()
}
