package main

import (
	"../items"
	"../wo"
	"time"
)

func spawnCoinsInWorld() {
	sleepTime := 10000 * time.Millisecond

	for {
		spawnCoin()
		time.Sleep(sleepTime)
		sleepTime += sleepTime
	}
}

func spawnCoin() {
	coord := wo.RandomEmptyCoord()
	coins := items.NewCoinStack()
	wo.AddItemsToStack(coord, coins)
}

func main() {
	spawnCoinsInWorld()
}
