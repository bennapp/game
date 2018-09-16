package rc

import (
	"../evt"
	"../gs"
	"../items"
	"../obj"
	"../os_util"
	"../pnt"
	"../store"
	"encoding/json"
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
	opt, err := redis.ParseURL(os_util.GetEnv("REDIS_URL", "redis://localhost:6379/0"))
	if err != nil {
		panic(err)
	}

	REDIS_INSTANCE = &RedisManager{
		client: redis.NewClient(opt),
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

func (manager *RedisManager) SaveItemsLocation(coord gs.Coord, items *items.Items) {
	itemsLocationStore := store.NewItemsLocationStore(coord, items)
	manager.set(itemsLocationStore)
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

	if objectData == "" {
		return nil
	}

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

	if objectStore == nil {
		return nil
	}

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

func (manager *RedisManager) LoadItemsStoreFromCoord(coord gs.Coord) *store.ItemsLocationStore {
	itemsStore := store.NewItemsStoreRetriever(coord)
	serializedString := manager.get(itemsStore)

	if serializedString == "" {
		return nil
	}

	itemsStore.SerializedItems = []byte(serializedString)

	return itemsStore
}

func (manager *RedisManager) SubscribeToObjectChannel(object obj.Objectable) chan string {
	pubsub := manager.client.Subscribe(object.ObjectId())

	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	redisChannel := pubsub.Channel()

	ch := make(chan string)
	go func() {
		for {
			select {
			case message := <-redisChannel:
				ch <- message.Payload
			default:
				// no op
			}
		}
	}()

	return ch
}

func (manager *RedisManager) WriteToObjectEventChannel(object obj.Objectable, event *evt.Event) {
	serializedEvent, _ := json.Marshal(event)
	manager.client.Publish(object.ObjectId(), serializedEvent)
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
