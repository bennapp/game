package wo

import (
	"../dbs"
	"../gs"
	"../obj"
)

func MoveObject(movable obj.Movable, vector gs.Vector) {
	coord := movable.GetLocation().AddVector(vector)
	cell := dbs.LoadCell(coord)

	if cell.IsMovableThrough() {
		dbs.DeleteObjectLocation(movable.GetLocation(), movable.(obj.Objectable))
		dbs.SaveObjectAndLocation(coord, movable.(obj.Objectable))
		movable.SetLocation(coord)
	}
}
