package main

import (
	"./wo"
	"fmt"
	"time"
)

func spawnCoinsInWorld() {
	sleepTime := 10000 * time.Millisecond

	for {
		fmt.Printf("basic.go: spawning random coin.\n")
		spawnCoin()
		time.Sleep(sleepTime)
		sleepTime += sleepTime
	}
}

func spawnCoin() {
	coin := wo.BuildAndStoreCoin()
	wo.PlaceCoinRandomly(coin)
}

func main() {
	fmt.Printf("coin.go: running World Elements\n")

	wo.Init()
	spawnCoinsInWorld()
}
