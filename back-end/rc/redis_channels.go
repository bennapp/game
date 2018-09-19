package rc

import (
	"../evt"
	"../gs"
	"../obj"
	"encoding/json"
	"github.com/go-redis/redis"
)

const WORLD_GENERATION_CHANNEL = "world_generation"

func (manager *RedisManager) SubscribeToWorldGenerationChannel() chan string {
	return manager.redisSubscribe(WORLD_GENERATION_CHANNEL)
}

func (manager *RedisManager) SubscribeToObjectChannel(object obj.Objectable) chan string {
	return manager.redisSubscribe(object.ObjectId())
}

func (manager *RedisManager) redisSubscribe(channelName string) chan string {
	pubsub := manager.client.Subscribe(channelName)

	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	redisChannel := pubsub.Channel()

	ch := make(chan string)
	go relayRedisPayload(ch, redisChannel)

	return ch
}

func relayRedisPayload(ch chan string, redisChannel <-chan *redis.Message) {
	for {
		message := <-redisChannel
		ch <- message.Payload
	}
}

func (manager *RedisManager) WriteToWorldGenerationChannel(coord gs.Coord) {
	manager.client.Publish(WORLD_GENERATION_CHANNEL, coord.Key())
}

func (manager *RedisManager) WriteToObjectEventChannel(object obj.Objectable, event *evt.Event) {
	serializedEvent, _ := json.Marshal(event)
	manager.client.Publish(object.ObjectId(), serializedEvent)
}
