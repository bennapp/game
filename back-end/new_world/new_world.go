package main

import (
	"../rc"
	"fmt"
)

func main() {
	fmt.Printf("newWorld.go: resetting world\n")
	rc.Manager().FlushAll()

	fmt.Printf("newWorld.go: creating new world\n")
	// zerozero := gs.NewCoord(0, 0)
	// wg.GenerateWorld(zerozero)
}
