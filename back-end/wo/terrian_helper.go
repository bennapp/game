package wo

import (
	"../dbs"
	"../gs"
	"../pnt"
	"../terr"
)

func SetTerrain(coord gs.Coord, terrain terr.Terrain) {
	paint := dbs.LoadPaintByCoord(coord)

	if paint == nil {
		paint = pnt.NewPaintByTerrain(terrain)
	} else {
		paint.Terrain = terrain
	}

	dbs.SavePaintLocation(coord, paint)
}
