package dbs

import (
	"../gs"
	"../obj"
	"../pnt"
	"../rc"
)

func SaveObject(object obj.Objectable) {
	rc.Manager().SaveObject(object)
}

func SaveObjectLocation(coord gs.Coord, object obj.Objectable) {
	rc.Manager().SaveObjectLocation(coord, object)
}

func SavePaintLocation(coord gs.Coord, paint *pnt.Paint) {
	rc.Manager().SavePaintLocation(coord, paint)
}
