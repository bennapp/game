package tests

import (
	"../dbs"
	"../gs"
	"../obj"
	"../pnt"
	"../terr"
	"testing"
)

func TestLoaderAndSaveCell(t *testing.T) {
	coord := gs.NewRandomCoord()

	coin := obj.NewCoin()
	paint := pnt.NewPaint(terr.NewRock(), pnt.Zone{})

	dbs.SaveObject(coin)
	dbs.SaveObjectLocation(coord, coin)
	dbs.SavePaintLocation(coord, paint)

	newCell := dbs.LoadCell(coord)
	newCoin := dbs.LoadObject(coin.ObjectId())

	t.Log("loaded cell from db:")
	t.Log(newCell)

	t.Log("loaded obj from db:")
	t.Log(newCoin)

	emptyCoord := gs.NewCoord(-1, -1)
	emptyCell := dbs.LoadCell(emptyCoord)

	t.Log(emptyCell)
	if emptyCell.IsMovableThrough() == false {
		t.Error("empty cell should be movable through")
	}

	if newCell.IsMovableThrough() == true {
		t.Error("cell with an object in it should not be movable through")
	}

	dbs.DeleteObjectLocation(coord, coin)
	cellWithoutCoin := dbs.LoadCell(coord)

	if cellWithoutCoin.Object != nil {
		t.Error("cell should have an empty object now that the coin has been removed")
	}

	dbs.DeleteObject(coin)
	emptyObject := dbs.LoadObject(coin.ObjectId())

	if emptyObject != nil {
		t.Error("Coin should be nil now that it is deleted")
	}
}
