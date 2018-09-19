package movs

import (
	"../dbs"
	"../gs"
	"../mov"
	"../obj"
	"../wg"
)

func MoveObject(movable mov.Movable, vector gs.Vector) {
	originalCoord := movable.GetLocation()

	nextCoord := originalCoord.AddVector(vector)
	nextCell := dbs.LoadCell(nextCoord)

	if nextCell.IsMovableThrough() {
		if obj.IsPlayer(movable) {
			distantCoord := nextCoord.AddVector(vector.Scale(gs.WORLD_GENERATION_DISTANCE / 2))
			wg.DetectWorldGeneration(distantCoord)
		}

		dbs.DeleteObjectLocation(originalCoord, movable.(obj.Objectable))
		movable.SetLocation(nextCoord)
		dbs.SaveObjectAndLocation(nextCoord, movable.(obj.Objectable))
	}
}
