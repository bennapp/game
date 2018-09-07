package gs

type Coord struct {
	X int
	Y int
}

func NewCoord(x int, y int) Coord {
	return Coord{X: x, Y: y}
}
