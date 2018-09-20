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

		regionStartX := (coord.X / gs.WORLD_GENERATION_DISTANCE) * gs.WORLD_GENERATION_DISTANCE
		regionStartY := (coord.Y / gs.WORLD_GENERATION_DISTANCE) * gs.WORLD_GENERATION_DISTANCE

		coordRegionStart := gs.NewCoord(regionStartX, regionStartY)

		go GenerateWorld(coordRegionStart)
	}
}
