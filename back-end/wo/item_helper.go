package wo

import (
	"../dbs"
	"../gs"
	"../items"
)

func AddItemsToStack(coord gs.Coord, stack items.ItemStack) {
	itemsToAdd := dbs.LoadItemsByCoord(coord)

	if itemsToAdd == nil {
		itemsToAdd = items.NewItemsWith(stack)
	} else {
		itemsToAdd.ItemStacks = append(itemsToAdd.ItemStacks, stack)
	}

	dbs.SaveItemsLocation(coord, itemsToAdd)
}
