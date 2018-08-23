package rc

import "github.com/go-redis/redis"

var INSTANCE RedisManager

type RedisManager struct {
	client *redis.Client
}

func Manager() RedisManager {
	if &INSTANCE == nil {
		initializeRedisClient()
	}

	return INSTANCE
}

func initializeRedisClient() {
	INSTANCE = RedisManager{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (manager *RedisManager) Client(v interface{}) *redis.Client {
	//TODO : return different redis client for different types
	return INSTANCE.client
}
