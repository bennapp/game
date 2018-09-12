package tests

import (
	"../dbs"
	"../gs"
	"../items"
	"../obj"
	"../pnt"
	"../terr"
	"testing"
)

func TestLoaderAndSaveCell(t *testing.T) {
	coord := gs.NewRandomCoord()

	player := obj.NewPlayer()
	paint := pnt.NewPaint(terr.NewRock(), pnt.Zone{})

	dbs.SaveObject(player)
	dbs.SaveObjectLocation(coord, player)
	dbs.SavePaintLocation(coord, paint)

	newCell := dbs.LoadCell(coord)
	newCoin := dbs.LoadObject(player.ObjectId())

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

	dbs.DeleteObjectLocation(coord, player)
	cellWithoutCoin := dbs.LoadCell(coord)

	if cellWithoutCoin.Object != nil {
		t.Error("cell should have an empty object now that the player has been removed")
	}

	dbs.DeleteObject(player)
	emptyObject := dbs.LoadObject(player.ObjectId())

	if emptyObject != nil {
		t.Error("Coin should be nil now that it is deleted")
	}
}

func TestSavingAndLoadingItems(t *testing.T) {
	coord := gs.NewRandomCoord()

	coinStack := items.NewItemStack(5, items.COIN)
	woodStack := items.NewItemStack(2, items.WOOD)
	itemStacks := []items.ItemStack{coinStack, woodStack}
	someItems := items.NewItems(itemStacks)

	dbs.SaveItemsLocation(coord, someItems)

	t.Log(dbs.LoadCell(coord).Items)
}
