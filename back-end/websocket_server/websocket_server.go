package main

import (
	"../rc"
	"../ws"
)

func main() {
	// This is for debugging and getting a clean state every server restart
	rc.Manager().FlushAll()

	ws.RunServer()
}
