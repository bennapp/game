package tests

import (
	"../dbs"
	"../gs"
	"../movs"
	"../obj"
	"../pnt"
	"../rc"
	"../terr"
	"testing"
	"time"
)

func TestMovement(t *testing.T) {
	defer func() {
		rc.Manager().FlushAll()
	}()

	coord := gs.NewCoord(5, 5)
	player := obj.NewPlayerAt(coord)
	dbs.SaveObjectAndLocation(coord, player)
	previousLocation := player.GetLocation()
	t.Log(previousLocation)

	movs.MoveObject(player, gs.NewVector(0, 1))

	location := player.GetLocation()
	if location.X != 5 {
		t.Error("X should not change")
	}

	if location.Y != 6 {
		t.Error("Y should increment")
	}

	cell := dbs.LoadCell(player.GetLocation())
	if cell.Object.ObjectId() != player.ObjectId() {
		t.Error("Cell should have player in it")
	}

	previousCell := dbs.LoadCell(previousLocation)
	if previousCell.Object != nil {
		t.Log(previousCell)
		t.Error("Previous cell should not have the player inside of it")
	}
}

func TestMovementRegulator(t *testing.T) {
	defer func() {
		rc.Manager().FlushAll()
	}()

	coord := gs.NewCoord(5, 5)
	player := obj.NewPlayerAt(coord)
	dbs.SaveObjectAndLocation(coord, player)

	mudCoord := gs.NewCoord(5, 6)
	mudPaint := pnt.NewPaint(terr.NewMud(), pnt.Zone{})
	dbs.SavePaintLocation(mudCoord, mudPaint)

	movs.RegulateMove(player, gs.NewVector(0, 1))
	movs.RegulateMove(player, gs.NewVector(0, 1))
	movs.RegulateMove(player, gs.NewVector(-1, 0))

	// prints original location + moving through three locations
	for i := 0; i < 5; i++ {
		t.Log(player.GetLocation())
		time.Sleep(200 * time.Millisecond)
	}
}
