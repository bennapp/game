package cell

import (
	"../items"
	"../obj"
	"../pnt"
	"fmt"
)

type Cell struct {
	Paint  *pnt.Paint
	Items  *items.Items
	Object obj.Objectable
}

func (cell *Cell) String() string {
	return fmt.Sprintf("paint: %v, obj: %v", cell.Paint, cell.Object)
}

func NewCell(paint *pnt.Paint, items *items.Items, object obj.Objectable) *Cell {
	return &Cell{Paint: paint, Items: items, Object: object}
}

func (cell *Cell) isEmptyObject() bool {
	return cell.Object == nil
}

func (cell *Cell) isEmptyPaint() bool {
	return cell.Paint == nil
}

func (cell *Cell) isPaintMovableThrough() bool {
	return cell.isEmptyPaint() || cell.Paint.Terrain.Permeable
}

func (cell *Cell) IsMovableThrough() bool {
	return cell.isEmptyObject() && cell.isPaintMovableThrough()
}
