package movs

import (
	"../dbs"
	"../gs"
	"../mov"
	"../obj"
)

func MoveObject(movable mov.Movable, vector gs.Vector) {
	coord := movable.GetLocation().AddVector(vector)
	cell := dbs.LoadCell(coord)

	if cell.IsMovableThrough() {
		dbs.DeleteObjectLocation(movable.GetLocation(), movable.(obj.Objectable))
		movable.SetLocation(coord)
		dbs.SaveObjectAndLocation(coord, movable.(obj.Objectable))
	}
}
