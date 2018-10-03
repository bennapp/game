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
			scaledVector := vector.Scale(gs.WORLD_GENERATION_DISTANCE / 2)
			distantCoord := originalCoord.AddVector(scaledVector)
			wg.DetectWorldGeneration(distantCoord)
		}

		dbs.DeleteObjectLocation(originalCoord, movable.(obj.Objectable))
		movable.SetLocation(nextCoord)
		dbs.SaveObjectAndLocation(nextCoord, movable.(obj.Objectable))
	}
}
