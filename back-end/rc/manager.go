package rc

import (
	"../gs"
	"github.com/go-redis/redis"
)

var REDIS_INSTANCE *RedisManager

type RedisManager struct {
	client *redis.Client
}

func Manager() *RedisManager {
	if REDIS_INSTANCE == nil {
		initializeRedisClient()
	}

	return REDIS_INSTANCE
}

func initializeRedisClient() {
	REDIS_INSTANCE = &RedisManager{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func (manager *RedisManager) FlushAll() {
	manager.client.FlushAll()
}

func (manager *RedisManager) SaveObjectLocation(coord gs.Coord, object Dbo) {
	objectLocationStore := newObjectLocationStore(coord, object)
	manager.set(objectLocationStore)
}

func (manager *RedisManager) SaveObject(object Dbo) {
	objectStore := newObjectStore(object)
	manager.set(objectStore)
}

func (manager *RedisManager) DeleteObject(object Dbo) {
	objectStore := newObjectStore(object)
	manager.delete(objectStore)
}

func (manager *RedisManager) LoadObjectFromCoord(coord gs.Coord) Dbo {
	objectLocationStore := newObjectLocationStoreRetriever(coord)

	objectId, _ := manager.get(objectLocationStore)
	// handle nil or empty

	newObjectStoreRetriever(objectId)
	// deserialize this
}

//func (manager *RedisManager) LoadFromKey(key string) string {
//
//}

func (manager *RedisManager) set(store RedisStore){
	manager.client.Set(store.Key(), store.Value(), 0).Err()
}

func (manager *RedisManager) delete(store RedisStore){
	manager.client.Set(store.Key(), nil, 0).Err()
}

func (manager *RedisManager) get(store RedisStore) string {
	val, _ := manager.client.Get(store.Key()).Result()
	return val
}
