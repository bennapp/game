package gs

import "fmt"

type Coord struct {
	X int
	Y int
}

func (coord Coord) Key() string {
	return fmt.Sprintf("%v,%v", coord.X, coord.Y)
}

func NewCoord(x int, y int) Coord {
	return Coord{X: x, Y: y}
}
