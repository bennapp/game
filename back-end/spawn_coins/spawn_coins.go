package main

import (
	"../items"
	"../wo"
	"time"
	"fmt"
)

func spawnCoinsInWorld() {
	sleepTime := 10000 * time.Millisecond

	for {
		fmt.Printf("spawning coin\n")
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
