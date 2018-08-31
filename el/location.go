package el

import (
	"../gs"
	"fmt"
)

type Location struct {
	GridCoord gs.Coord
}

func (location *Location) LocationKey() string {
	return fmt.Sprintf("%v", location.GridCoord.Key())
}

func NewLocation(gridCoord gs.Coord) *Location {
	return &Location{GridCoord: gridCoord}
}
