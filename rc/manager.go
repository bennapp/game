package rc

import (
	"github.com/go-redis/redis"
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

func initializeRedisClient() {
	INSTANCE = &RedisManager{
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

func (manager *RedisManager) Save(o Dbo) {
	client := manager.Client(o.Key())
	//fmt.Printf("manager.go: Client found for key: %s\n", o.Key())

	serialize := o.Serialize()
	//fmt.Printf("manager.go: Serialize success: %s\n", serialize)

	err := client.Set(o.Key(), serialize, 0).Err()
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
	//fmt.Printf("manager.go: Loading from Key: %s\n", key)

	val, _ := manager.Client(key).Get(key).Result()

	//if err != nil {
	//	panic(err)
	//}

	return val
}
