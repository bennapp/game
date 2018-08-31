package wo

import "math/rand"

import (
	"../gs"
)

func randomPair(n int) (int, int) {
	return rand.Intn(n), rand.Intn(n)
}

func RandomSubWorldCoord() gs.Coord {
	x, y := randomPair(gs.WORLD_SIZE)

	return gs.NewCoord(x, y)
}

func isOutOfBound(x int, y int, bound int) bool {
	return x < 0 || y < 0 || x >= bound || y >= bound
}

func wrap(base int, add int, max int) int {
	sum := base + add

	if sum > 0 {
		return sum % max
	} else {
		return ((sum % max) + max) % max
	}
}

func carry(base int, add int, max int) int {
	sum := base + add

	if sum > 0 {
		return sum / max
	} else {
		return (sum - max + 1) / max
	}
}
