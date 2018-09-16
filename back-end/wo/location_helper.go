package wo

import (
	"../dbs"
	"../gs"
)

func RandomEmptyCoord() gs.Coord {
	var coord gs.Coord

	for {
		coord = gs.NewRandomCoord()
		cell := dbs.LoadCell(coord)
		if cell.IsMovableThrough() {
			break
		}
	}

	return coord
}
