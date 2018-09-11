package cell_test

import (
	"../cell"
	"../items"
	"../obj"
	"../pnt"
	"../terr"
	"testing"
)

func TestCells(t *testing.T) {
	cell := cell.NewCell(
		pnt.NewPaint(terr.NewRock(), pnt.Zone{}),
		items.NewItems(make([]items.ItemStack, 0)),
		obj.NewCoin(),
	)
	t.Log(cell)
}
