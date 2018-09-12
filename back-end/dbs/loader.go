package dbs

import (
	"../cell"
	"../gs"
	"../items"
	"../obj"
	"../pnt"
	"../rc"
	"../typf"
)

func LoadObject(objectId string) obj.Objectable {
	objectStore := rc.Manager().LoadObjectStore(objectId)

	if objectStore == nil {
		return nil
	}

	object := typf.Factory().DeserializeObject(objectStore)

	return object.(obj.Objectable)
}

func LoadObjectByCoord(coord gs.Coord) obj.Objectable {
	objectStore := rc.Manager().LoadObjectStoreFromCoord(coord)

	var object obj.Objectable
	if objectStore != nil {
		object = typf.Factory().DeserializeObject(objectStore)
	}
	return object
}

func LoadCell(coord gs.Coord) *cell.Cell {
	object := LoadObjectByCoord(coord)

	paintStore := rc.Manager().LoadPaintStoreFromCoord(coord)

	var paint *pnt.Paint
	if paintStore != nil {
		paint = typf.Factory().DeserializePaint(paintStore)
	}

	itemsStore := rc.Manager().LoadItemsStoreFromCoord(coord)

	var items *items.Items
	if itemsStore != nil {
		items = typf.Factory().DeserializeItems(itemsStore)
	}

	return cell.NewCell(paint, items, object)
}
