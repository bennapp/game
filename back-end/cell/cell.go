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

func (cell *Cell) IsEmpty() bool {
	return cell.isEmptyObject() && cell.isEmptyPaint() && cell.isEmptyItems()
}

func (cell *Cell) isEmptyObject() bool {
	return cell.Object == nil
}

func (cell *Cell) isEmptyPaint() bool {
	return cell.Paint == nil
}

func (cell *Cell) isEmptyItems() bool {
	return cell.Items == nil
}

func (cell *Cell) isPaintMovableThrough() bool {
	return cell.isEmptyPaint() || cell.Paint.Terrain.Permeable
}

func (cell *Cell) IsMovableThrough() bool {
	return cell.isEmptyObject() && cell.isPaintMovableThrough()
}

func (cell *Cell) DebugString() string {
	if cell.IsEmpty() {
		return " "
	}

	if !cell.isEmptyObject() {
		return "O"
	}

	if !cell.isEmptyPaint() {
		return "P"
	}

	if !cell.isEmptyItems() {
		return "I"
	}

	return cell.String()
}
