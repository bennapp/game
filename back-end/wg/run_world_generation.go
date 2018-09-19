package wg

import (
	"../gs"
	"../rc"
)

func RunWorldGeneration() {
	worldGenerationChannel := rc.Manager().SubscribeToWorldGenerationChannel()

	for {
		coordKey := <-worldGenerationChannel
		coord := gs.NewCoordFromKey(coordKey)
		go generateWorld(coord)
	}
}
