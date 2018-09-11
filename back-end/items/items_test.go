package items_test

import (
	"../items"
	"testing"
)

func TestItems(t *testing.T) {
	coinStack := items.NewItemStack(5, items.COIN)
	woodStack := items.NewItemStack(2, items.WOOD)
	itemStacks := []items.ItemStack{coinStack, woodStack}
	someItems := items.NewItems(itemStacks)

	t.Log(someItems)
}
