package items

import (
	"../typ"
)

const ITEMS = "items"

type Items struct {
	typ.Type
	ItemStacks []ItemStack
}

type ItemStack struct {
	Amount int
	Item
}

type Item struct {
	ItemType string
}

type Itemable interface {
	GetItemType() string
}

func (item *Item) GetItemType() string {
	return item.ItemType
}

func NewItemsWith(stack ItemStack) *Items {
	itemStacks := []ItemStack{stack}
	return NewItems(itemStacks)
}

func NewItems(itemStacks []ItemStack) *Items {
	return &Items{Type: typ.NewType(ITEMS), ItemStacks: itemStacks}
}

func LoadItems() typ.Typical {
	return &Items{}
}

func NewItemStack(amount int, itemType string) ItemStack {
	return ItemStack{Amount: amount, Item: newItem(itemType)}
}

func newItem(itemType string) Item {
	return Item{ItemType: itemType}
}
