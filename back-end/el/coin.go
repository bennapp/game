package el

import (
	"github.com/google/uuid"
)

const COIN = "coin"

type Coin struct {
	Amount int
	Type   string
	Object
}

func (coin *Coin) String() string {
	return "C"
}

func (coin *Coin) Id() string {
	return coin.Object.Id.String()
}

func newCoin() *Coin {
	uuid, _ := uuid.NewUUID()
	return &Coin{Type: COIN, Object: Object{Id: uuid}}
}
