package dbs

import (
	"../gs"
	"../items"
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

func SaveItemsLocation(coord gs.Coord, items *items.Items) {
	rc.Manager().SaveItemsLocation(coord, items)
}
