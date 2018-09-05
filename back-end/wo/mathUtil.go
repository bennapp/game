package wo

import "math/rand"

func randomPair(n int) (int, int) {
	return rand.Intn(n), rand.Intn(n)
}

func isOutOfBound(x int, y int, bound int) bool {
	return x < 0 || y < 0 || x >= bound || y >= bound
}
