package tests

import (
	"../gs"
	"../rc"
	"github.com/go-redis/redis"
	"testing"
)

func TestEvents(t *testing.T) {
	eventChannel := make(chan string)

	ch := rc.Manager().SubscribeToCoordEvents(coord)

	go propagateCoordEvent(ch, eventChannel)

	rc.Manager().WriteToCoordEvents(coord, "hello")
}

func propagateCoordEvent(ch <-chan *redis.Message, eventChannel chan string) {
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				return
			}
			eventChannel <- message.String()
		default:
		}
	}
}
