package main

import (
	"./wo"
	"fmt"
)

func main() {
	fmt.Printf("newWorld.go: resetting world\n")
	wo.ResetWorld()

	fmt.Printf("newWorld.go: creating new world\n")
	wo.InitializeWorld()
}
