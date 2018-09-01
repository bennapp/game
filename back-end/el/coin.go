package el

import (
	"../rc"
)

const COIN = "coin"

type Coin struct {
	Amount int
	id     int
	Type   string
}

func (coin *Coin) String() string {
	return "C"
}

func (coin *Coin) Id() int {
	return coin.id
}

func (coin *Coin) Key() string {
	return rc.GenerateKey(COIN, coin.id)
}

func newCoinDbo(id int) rc.Dbo {
	return &Coin{id: id, Type: COIN}
}
