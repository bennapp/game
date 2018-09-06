package el

import (
	"../rc"
	"github.com/google/uuid"
)

const COIN = "coin"

type Coin struct {
	Amount int
	Id     uuid.UUID
	Type   string
}

func (coin *Coin) String() string {
	return "C"
}

func (coin *Coin) Key() string {
	return rc.GenerateKey(coin.Id)
}

func newCoinDbo(id uuid.UUID) rc.Dbo {
	return &Coin{Id: id, Type: COIN}
}

func (coin *Coin) Load() rc.Dbo {
	return &Coin{Type: COIN}
}
