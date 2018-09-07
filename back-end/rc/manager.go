package rc

import (
	"../gs"
	"fmt"
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

func (manager *RedisManager) PrintAllKeyValuesForDebugging() {
	results, _ := manager.client.Keys("*").Result()

	for _, str := range results {
		value, _ := manager.client.Get(str).Result()
		fmt.Printf("key: %v, value: %v\n", str, value)
	}
}

func (manager *RedisManager) FlushAll() {
	manager.client.FlushAll()
}

func (manager *RedisManager) SaveObjectLocation(coord gs.Coord, object Dbo) {
	objectLocationStore := newObjectLocationStore(coord, object)
	manager.set(objectLocationStore)
}

func (manager *RedisManager) DeleteObjectLocation(coord gs.Coord, object Dbo) {
	objectLocationStore := newObjectLocationStore(coord, object)
	manager.delete(objectLocationStore)
}

func (manager *RedisManager) SaveObject(object Dbo) {
	objectStore := newObjectStore(object)
	manager.set(objectStore)
}

func (manager *RedisManager) DeleteObject(object Dbo) {
	objectStore := newObjectStore(object)
	manager.delete(objectStore)
}

func (manager *RedisManager) LoadObjectTypeFromCoord(coord gs.Coord) (string, *ObjectStore) {
	objectLocationStore := newObjectLocationStoreRetriever(coord)

	objectId := manager.get(objectLocationStore)

	if objectId == "" {
		return "", nil
	}

	objectStore := newObjectStoreRetriever(objectId)
	objectStore.SerializedObject = []byte(manager.get(objectStore))

	objectType := newTypeDeserializer(objectStore.SerializedObject).Type

	return objectType, objectStore
}

func (manager *RedisManager) set(store RedisStore) {
	manager.client.Set(store.Key(), store.Value(), 0).Err()
}

func (manager *RedisManager) delete(store RedisStore) {
	manager.client.Set(store.Key(), nil, 0).Err()
}

func (manager *RedisManager) get(store RedisStore) string {
	val, _ := manager.client.Get(store.Key()).Result()
	return val
}
