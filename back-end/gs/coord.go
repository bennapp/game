package gs

import (
	"../math_util"
	"fmt"
	"strconv"
	"strings"
)

type Coord struct {
	X int
	Y int
}

func NewCoord(x int, y int) Coord {
	return Coord{X: x, Y: y}
}

func NewRandomCoord() Coord {
	x, y := math_util.RandomPair(RANDOM_CORD_SIZE)
	return Coord{X: x, Y: y}
}

func NewCoordFromKey(key string) Coord {
	positions := strings.Split(key, ",")
	x, _ := strconv.Atoi(positions[0])
	y, _ := strconv.Atoi(positions[1])

	return NewCoord(x, y)
}

func (coord Coord) Key() string {
	return fmt.Sprintf("%v,%v", coord.X, coord.Y)
}

func (coord Coord) AddVector(vector Vector) Coord {
	return NewCoord(coord.X+vector.X, coord.Y+vector.Y)
}
