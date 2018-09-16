package new_world

import (
	"../rc"
	"../wo"
	"fmt"
)

func main() {
	fmt.Printf("newWorld.go: resetting world\n")
	rc.Manager().FlushAll()

	fmt.Printf("newWorld.go: creating new world\n")
	wo.InitializeWorld()
}
