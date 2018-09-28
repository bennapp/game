package store

import (
	"../gs"
	"fmt"
)

const WORLD_GENERATION_LAYER = "g"
const WORLD_GENERATION_EXPLORED_VALUE = "t"

type WorldGenerationStore struct {
	Coord      gs.Coord
	IsExplored bool
}

func (wgs *WorldGenerationStore) Key() string {
	coordRegion := coordRegion(wgs.Coord)

	return fmt.Sprintf("%v:%v", WORLD_GENERATION_LAYER, coordRegion.Key())
}

func (wgs *WorldGenerationStore) Value() string {
	return WORLD_GENERATION_EXPLORED_VALUE
}

func coordRegion(coord gs.Coord) gs.Coord {
	x := regionConvert(coord.X)
	y := regionConvert(coord.Y)

	return gs.NewCoord(x, y)
}

func regionConvert(position int) int {
	return position / gs.WORLD_GENERATION_DISTANCE
}

func NewWorldGenerationStore(coord gs.Coord) *WorldGenerationStore {
	return &WorldGenerationStore{Coord: coord}
}

func (wgs *WorldGenerationStore) Retrieve(storeData string) {
	wgs.IsExplored = storeData == WORLD_GENERATION_EXPLORED_VALUE
}
