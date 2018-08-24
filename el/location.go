package el

import (
	"../gs"
	"fmt"
)

type Location struct {
	SubWorldCoord gs.Coord
	GridCoord     gs.Coord
}

func (location *Location) LocationKey() string {
	return fmt.Sprintf("%v:%v", location.SubWorldCoord.Key(), location.GridCoord.Key())
}

func NewLocation(subWorldCoord gs.Coord, gridCoord gs.Coord) *Location {
	return &Location{SubWorldCoord: subWorldCoord, GridCoord: gridCoord}
}
