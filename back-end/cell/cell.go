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

// Right now this implies that you cannot move through a cell if an object is there. Which would mean we cannot move through
// coins, which is not 'correct' for our the engine right now. They way we could change / fix this is introduce an items location layer hmm...
func (cell *Cell) IsMovableThrough() bool {
	return cell.isEmptyObject() && cell.isPaintMovableThrough()
}
