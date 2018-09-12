package rc

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"os"
)

var INSTANCE *RedisManager

type RedisManager struct {
	client *redis.Client
}

func Manager() *RedisManager {
	if INSTANCE == nil {
		initializeRedisClient()
	}

	return INSTANCE
}

func GetEnv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func initializeRedisClient() {
	opt, err := redis.ParseURL(GetEnv("REDIS_URL", "redis://localhost:6379/0"))
	if err != nil {
		panic(err)
	}

	INSTANCE = &RedisManager{
		client: redis.NewClient(opt),
	}
}

func (manager *RedisManager) FlushAll() {
	//ALL Client flush
	manager.client.FlushAll()
}

func (manager *RedisManager) Client(v interface{}) *redis.Client {
	//TODO : return different redis client for different types
	return manager.client
}

func (manager *RedisManager) Save(o Dbo) {
	client := manager.Client(o.Key())

	serialize, err := json.Marshal(o)

	err = client.Set(o.Key(), serialize, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (manager *RedisManager) Delete(o Dbo) {
	err := manager.Client(o.Key()).Set(o.Key(), nil, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (manager *RedisManager) LoadFromId(t string, id int) string {
	return manager.LoadFromKey(GenerateKey(t, id))
}

func (manager *RedisManager) LoadFromKey(key string) string {
	val, _ := manager.Client(key).Get(key).Result()

	return val
}
