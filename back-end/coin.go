package main

import (
	"./gs"
	"./wo"
	"fmt"
	"time"
)

func spawnCoinsInWorld() {
	sleepTime := 10000 * time.Millisecond

	for {
		fmt.Printf("basic.go: spawning random coin.\n")
		randomSubWorldCoord := wo.RandomSubWorldCoord()
		spawnCoinInSubWorld(randomSubWorldCoord)
		time.Sleep(sleepTime)
		sleepTime += sleepTime
	}
}

func spawnCoinInSubWorld(subWorldCoord gs.Coord) {
	coin := wo.BuildAndStoreCoin()
	wo.PlaceCoin(subWorldCoord, coin)
}

func main() {
	fmt.Printf("coin.go: running World Elements\n")

	wo.Init()
	spawnCoinsInWorld()
}
