package tests

import (
	"../el"
	"../gs"
	"../rc"
	"fmt"
)

func DebugRedisStore() {
	fmt.Println("Printing all key values")
	rc.Manager().PrintAllKeyValuesForDebugging()

	fmt.Println()
	fmt.Println("Printing all loaded dbos")
	printAllDbos()
}

func printAllDbos() {
	for i := 0; i < gs.WORLD_SIZE; i++ {
		for j := 0; j < gs.WORLD_SIZE; j++ {
			coord := gs.NewCoord(i, j)

			dbo := el.Factory().LoadObjectFromCoord(coord)
			fmt.Printf("%+v\n", dbo)
		}
	}
}
