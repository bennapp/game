package math_util

import "math/rand"

func Abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func Max(a, b int) int {
	if a < b {
		return b
	}

	return a
}

func RandomPair(n int) (int, int) {
	return rand.Intn(n), rand.Intn(n)
}

func isOutOfBound(x int, y int, bound int) bool {
	return x < 0 || y < 0 || x >= bound || y >= bound
}
