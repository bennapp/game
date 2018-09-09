package dbs

import (
	"../gs"
	"../obj"
	"../rc"
)

func DeleteObjectLocation(coord gs.Coord, object obj.Objectable) {
	rc.Manager().DeleteObjectLocation(coord, object)
}

func DeleteObject(object obj.Objectable) {
	rc.Manager().DeleteObject(object)
}
