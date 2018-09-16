package wo

import (
	"../dbs"
	"../obj"
)

func CreatePlayer() *obj.Player {
	coord := RandomEmptyCoord()
	player := obj.NewPlayerAt(coord)
	dbs.SaveObjectAndLocation(coord, player)

	return player
}

func DeletePlayer(player *obj.Player) {
	dbs.DeleteObjectAtLocation(player.GetLocation(), player)
}
