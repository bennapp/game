package cell_test

import (
	"../cell"
	"../obj"
	"../pnt"
	"../terr"
	"testing"
)

func TestCells(t *testing.T) {
	cell := cell.NewCell(
		pnt.NewPaint(terr.NewRock(), pnt.Zone{}),
		obj.NewCoin(),
	)
	t.Log(cell)
}
