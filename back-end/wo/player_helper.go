package wo

import (
	"../dbs"
	"../gs"
	"../obj"
	"fmt"
)

func CreatePlayer() *obj.Player {
	coord := gs.NewCoord(0, 0)
	player := obj.NewPlayerAt(coord)
	dbs.SaveObjectAndLocation(coord, player)

	return player
}

func DeletePlayer(player *obj.Player) {
	fmt.Println("deleting player?")
	dbs.DeleteObjectAtLocation(player.GetLocation(), player)
}
