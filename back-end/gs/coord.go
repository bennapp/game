package gs

import (
	"../util"
)

type Coord struct {
	X int
	Y int
}

func NewCoord(x int, y int) Coord {
	return Coord{X: x, Y: y}
}

func NewRandomCoord() Coord {
	x, y := util.RandomPair(WORLD_SIZE)
	return Coord{X: x, Y: y}
}

func (coord Coord) AddVector(vector Vector) Coord {
	return NewCoord(coord.X+vector.X, coord.Y+vector.Y)
}
