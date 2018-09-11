package dbs

import (
	"../cell"
	"../el"
	"../gs"
	"../items"
	"../obj"
	"../pnt"
	"../rc"
)

func LoadObject(objectId string) obj.Objectable {
	objectStore := rc.Manager().LoadObjectStore(objectId)

	if objectStore == nil {
		return nil
	}

	object := el.Factory().DeserializeObject(objectStore)

	return object.(obj.Objectable)
}

func LoadCell(coord gs.Coord) *cell.Cell {
	objectStore := rc.Manager().LoadObjectStoreFromCoord(coord)

	var object obj.Objectable
	if objectStore != nil {
		object = el.Factory().DeserializeObject(objectStore)
	}

	paintStore := rc.Manager().LoadPaintStoreFromCoord(coord)

	var paint *pnt.Paint
	if paintStore != nil {
		paint = el.Factory().DeserializePaint(paintStore)
	}

	itemsStore := rc.Manager().LoadItemsStoreFromCoord(coord)

	var items *items.Items
	if itemsStore != nil {
		items = el.Factory().DeserializeItems(itemsStore)
	}

	return cell.NewCell(paint, items, object)
}
