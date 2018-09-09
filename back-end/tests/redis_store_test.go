package tests

import (
	"../el"
	"../gs"
	"../rc"
	"testing"
)

func TestRedisStore(t *testing.T) {
	t.Log("Printing all key values")
	rc.Manager().PrintAllKeyValuesForDebugging()

	t.Log()
	t.Log("Printing all loaded objects and paint")
	printAllObjects(t)
}

func printAllObjects(t *testing.T) {
	for i := 0; i < gs.WORLD_SIZE; i++ {
		for j := 0; j < gs.WORLD_SIZE; j++ {
			coord := gs.NewCoord(i, j)

			object := el.Factory().LoadObjectFromCoord(coord)
			t.Logf("object: %v\n", object)

			paint := el.Factory().LoadPaintFromCoord(coord)
			t.Logf("paint: %v\n", paint)
		}
	}
}
