package obj

import (
	"math/rand"
)

const COIN = "coin"
const MAX_COIN_AMOUNT = 10

type Coin struct {
	Object
	Amount int
}

//func (coin *Coin) String() string {
//	return "C"
//}

func NewCoin() *Coin {
	amount := rand.Intn(MAX_COIN_AMOUNT)
	return &Coin{Object: newObject(COIN), Amount: amount}
}

func LoadCoin() Objectable {
	return &Coin{Object: loadObject(COIN)}
}
