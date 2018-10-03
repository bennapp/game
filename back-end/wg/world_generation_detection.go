package wg

import (
	"../dbs"
	"../gs"
	"../rc"
)

func DetectWorldGeneration(coord gs.Coord) {
	if !dbs.RegionExplored(coord) {
		dbs.SetRegionExplored(coord)
		rc.Manager().WriteToWorldGenerationChannel(coord)
	}
}
