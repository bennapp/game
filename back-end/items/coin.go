package items

import "math/rand"

const COIN = "coin"
const MAX_COIN_AMOUNT = 10

func NewCoin() Item {
	return newItem(COIN)
}

func NewCoinStack() ItemStack {
	amount := rand.Intn(MAX_COIN_AMOUNT)
	return ItemStack{Item: NewCoin(), Amount: amount}
}
