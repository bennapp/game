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

func newCoinDbo(isNew bool) rc.Dbo {
	coin := &Coin{Type: COIN}

	if isNew {
		coin.Id = uuid.New()
	}

	return coin
}
