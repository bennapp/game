package main

import (
	"./wo"
	"fmt"
)

func main() {
	wo.Init()

	fmt.Printf("newWorld.go: resetting world\n")
	wo.Reset()

	fmt.Printf("newWorld.go: start initializing world\n")
	wo.InitializeWorld()
	fmt.Printf("newWorld.go: finish initializing world\n")
}
