package main

import (
	"./wo"
	"fmt"
)

func main() {
	wo.Init()
	player := wo.InitializePlayer()

	fmt.Printf("Created new player\nLocation: %s\nKey:%s\n",
		player.LocationKey(),
		player.Key(),
	)
	wo.Close()
}