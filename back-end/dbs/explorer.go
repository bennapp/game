package dbs

import (
	"../gs"
	"../rc"
)

func RegionExplored(coord gs.Coord) bool {
	worldGenerationStore := rc.Manager().LoadWorldGenerationStore(coord)
	return worldGenerationStore.IsExplored
}

func SetRegionExplored(coord gs.Coord) {
	rc.Manager().SaveWorldGenerationStore(coord)
}
