package tests

import (
	"../gs"
	"../rc"
	"testing"
	"github.com/go-redis/redis"
)

func TestEvents(t *testing.T) {
	coord := gs.NewCoord(5, 5)
	ch := rc.Manager().SubscribeToCoordEvents(coord)

	go listenToCoordEvent(ch)

	rc.Manager().WriteToCoordEvents(coord, "hello")
}

func listenToCoordEvent(ch <-chan *redis.Message) {

}
