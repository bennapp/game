package rc

import (
	"../gs"
	"../obj"
	"../pnt"
	"../store"
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

func (manager *RedisManager) SaveObjectLocation(coord gs.Coord, object obj.Objectable) {
	objectLocationStore := store.NewObjectLocationStore(coord, object)
	manager.set(objectLocationStore)
}

func (manager *RedisManager) SavePaintLocation(coord gs.Coord, paint *pnt.Paint) {
	paintLocationStore := store.NewPaintLocationStore(coord, paint)
	manager.set(paintLocationStore)
}

func (manager *RedisManager) DeleteObjectLocation(coord gs.Coord, object obj.Objectable) {
	objectLocationStore := store.NewObjectLocationStore(coord, object)
	manager.delete(objectLocationStore)
}

func (manager *RedisManager) SaveObject(object obj.Objectable) {
	objectStore := store.NewObjectStore(object)
	manager.set(objectStore)
}

func (manager *RedisManager) DeleteObject(object obj.Objectable) {
	objectStore := store.NewObjectStore(object)
	manager.delete(objectStore)
}

func (manager *RedisManager) LoadObjectStore(objectId string) *store.ObjectStore {
	objectStore := store.NewObjectStoreRetriever(objectId)
	objectData := manager.get(objectStore)
	objectStore.Retrieve(objectData)

	return objectStore
}

func (manager *RedisManager) LoadObjectStoreFromCoord(coord gs.Coord) *store.ObjectStore {
	objectLocationStore := store.NewObjectLocationStoreRetriever(coord)
	objectId := manager.get(objectLocationStore)

	if objectId == "" {
		return nil
	}

	objectStore := manager.LoadObjectStore(objectId)

	return objectStore
}

func (manager *RedisManager) LoadPaintStoreFromCoord(coord gs.Coord) *store.PaintLocationStore {
	paintStore := store.NewPaintStoreRetriever(coord)
	serializedString := manager.get(paintStore)

	if serializedString == "" {
		return nil
	}

	paintStore.SerializedPaint = []byte(serializedString)

	return paintStore
}

func (manager *RedisManager) set(store store.RedisStore) {
	manager.client.Set(store.Key(), store.Value(), 0).Err()
}

func (manager *RedisManager) delete(store store.RedisStore) {
	manager.client.Set(store.Key(), nil, 0).Err()
}

func (manager *RedisManager) get(store store.RedisStore) string {
	val, _ := manager.client.Get(store.Key()).Result()
	return val
}
